package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	fsc "github.com/heron-sense/gadk/flow-state-code"
	"github.com/heron-sense/gadk/logger"
	"github.com/heron-sense/gadk/subroutine"
	"github.com/heron-sense/gadk/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/tinylib/msgp/msgp"
	"io"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

type Engine interface {
	Run(ctx context.Context, listener net.Listener) error
	Stop() error
	RegisterRoute(directive string, subroutine subroutine.Subroutine) error
}

func NewEngine() *_Engine {
	//load balance strategy
	//read timeout
	//idle timeout
	//write timeout

	engine := gin.New()
	tracer, closer := tracing.InitJaeger("dispatch")

	opentracing.SetGlobalTracer(tracer)

	return &_Engine{
		engine:     engine,
		innerGroup: engine.Group("/inner"),
		outerGroup: engine.Group("/"),
		Tracing: struct {
			Tracer opentracing.Tracer
			io.Closer
		}{Tracer: tracer, Closer: closer},
	}
}

type _Engine struct {
	engine     *gin.Engine
	innerGroup *gin.RouterGroup
	outerGroup *gin.RouterGroup
	Tracing    struct {
		Tracer opentracing.Tracer
		io.Closer
	}
}

func (e *_Engine) Stop() error {

	// WaitShutdownSignals 等待接收SIGQUIT或者SIGINT信号，用于Main方法的wait参数时，可以通过Ctrl+C来停止进程
	// 使用SIGQUIT之外的信号会通过优雅退出的方式结束

	return nil
}

func (e *_Engine) Run(ctx context.Context, listener net.Listener) error {
	e.RegisterRoute("/log/stats", &subroutine.GetInfo{})

	var err error
	go func() {
		err = e.engine.RunListener(listener)
	}()
	if err != nil {
		return err
	}
	select {
	case <-ctx.Done():
		return e.Stop()
	}
}

func getConfig(path string) fsc.FlowStateCode {
	fileList, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Printf("%s", err)
		return fsc.FlowPermissionDenied
	}

	for _, file := range fileList {
		if file.IsDir() {
			continue
		}

		content, err := ioutil.ReadFile(path + string([]byte{os.PathSeparator}) + file.Name())
		if err != nil {
			logger.Error("read err:%s", err)
			return fsc.FlowPermissionDenied
		}
		content = content
	}

	return fsc.FlowFinished
}

func GetExtension() http.Header {
	hdr := http.Header{}

	return hdr
}

func (app *_Engine) RegisterRoute(directive string, subroutine subroutine.Subroutine) error {
	app.engine.Handle(http.MethodGet, directive, func(c *gin.Context) {
		// TODO get config

		fmt.Printf("what?:%+v", c.GetHeader("content-type"))

		spanCtx, _ := app.Tracing.Tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(GetExtension()))
		span := app.Tracing.Tracer.StartSpan(directive, ext.RPCServerOption(spanCtx))
		defer span.Finish()
		span.SetTag(tracing.TagTracingID, "tracing-id")
		span.SetTag(tracing.TagIslandTerm, "unknown")
		span.SetTag(tracing.TagIslandName, "island-name")

		urlPath := "/url/path"

		//span, _ := opentracing.StartSpanFromContext(c, "dispatch")
		//defer span.Finish()

		conv := make(opentracing.HTTPHeadersCarrier)
		ext.SpanKindRPCClient.Set(span)
		ext.HTTPUrl.Set(span, urlPath)
		ext.HTTPMethod.Set(span, "GET")

		err := span.Tracer().Inject(
			span.Context(),
			opentracing.HTTPHeaders,
			conv,
		)
		if err != nil {
			logger.Vital("inject err:%s", err)
		}

		convBuf := make([]byte, 0, 200)
		_ = conv.ForeachKey(func(key, val string) error {
			convBuf = append(convBuf, fmt.Sprintf("&%s=%s", key, val)...)
			return nil
		})

		logger.Vital("conv:ext=%s", string(convBuf))

		parseReq(c, subroutine)

		remainingMs := time.Duration(math.MaxUint16)
		ctx, _ := context.WithDeadline(c, time.Now().Add(time.Millisecond*remainingMs))
		rsp, err := subroutine.Handle(ctx)
		if err != nil {
			//error process
			return
		}

		data, _ := json.Marshal(rsp)

		//data, err := GenPayload(0, rsp)
		//if err != nil {
		//	//error process
		//	return
		//}

		c.Data(200, "application/json", data)
	})
	return nil

}

func parseReq(ctx *gin.Context, req msgp.Unmarshaler) error {
	const dataMaxLength = 0xFFFFFF

	data, err := ctx.GetRawData()
	if err != nil {
		//error process
		return err
	}

	err = json.Unmarshal(data, req)
	//remaining, err := req.UnmarshalMsg(data)
	if err != nil {
		//error process
		return err
	}

	span, _ := opentracing.StartSpanFromContext(ctx, "parse")
	defer span.Finish()
	return nil
}

func GenPayload(fsCode fsc.FlowStateCode, reply subroutine.Reply) ([]byte, error) {
	buf := make([]byte, 0, reply.Msgsize()+200)
	assembler := &strings.Builder{}
	assembler.WriteString(fmt.Sprintf("{\"state_code\":%d,\"data\":", fsCode))

	buf, err := reply.MarshalMsg(buf)
	if err != nil {
		logger.Error("service unavailable:%s", err)
		return nil, err
	}

	assembler.WriteByte('}')
	return buf, nil
}
