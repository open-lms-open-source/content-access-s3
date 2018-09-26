package s3wrapper

import (
	"context"
	"errors"
	"net/http/httptrace"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/s3"

	"moodle-backup-filler/logger"
)

// S3 provides TTFB/Retry functionality for S3 calls.
type S3 struct {
	*s3.S3
}

// New creates a new instance of the S3 client wrapper.
func New(p client.ConfigProvider, cfgs *aws.Config) *S3 {
	return &S3{
		S3: s3.New(p, cfgs),
	}
}

// ttfbTime is a simple struct for our timeout context to keep track of the start time with.
type ttfbTime struct {
	time         time.Time
	gotFirstByte bool
	cancel       *context.CancelFunc
}

func (t *ttfbTime) cancelConnection() {
	if !t.gotFirstByte {
		logger.Err.Debugf("S3 source: TTFB timelimit exceeded: %d", time.Since(t.time)/time.Millisecond)
		(*t.cancel)()
	}
}

// getTTFBTimeoutContext provides a context with a TTFB timeout attached.
func getTTFBTimeoutContext(ttfbTimeout int64) context.Context {
	start := &ttfbTime{
		time:         time.Now(),
		gotFirstByte: false,
	}

	trace := &httptrace.ClientTrace{
		GotFirstResponseByte: func() {
			start.gotFirstByte = true
			logger.Err.Debugf("S3 source: TTFB: %d", time.Since(start.time)/time.Millisecond)
		},
	}

	ctx, cancel := context.WithCancel(httptrace.WithClientTrace(context.Background(), trace))
	start.cancel = &cancel

	time.AfterFunc(time.Duration(ttfbTimeout)*time.Millisecond, start.cancelConnection)

	return ctx
}

// GetObjectWithRetry will cancel its request after the provided timeout and retry exactly once.
func (s3 *S3) GetObjectWithRetry(input *s3.GetObjectInput, ttfbTimeout int64, ttfbRetries int64) (resp *s3.GetObjectOutput, err error) {
	var attempt int64 = 1

	delay := ttfbTimeout
	for attempt <= ttfbRetries {
		ctx := getTTFBTimeoutContext(delay)
		resp, err = s3.GetObjectWithContext(ctx, input)

		if ctx.Err() != context.Canceled {
			return resp, err
		}

		attempt++
		delay *= 2
	}

	err = errors.New("Request to S3 timed out")

	return resp, err
}
