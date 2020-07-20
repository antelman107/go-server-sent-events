package main

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Sender struct {
	Chan chan MessageInterface
}

var (
	ErrWriterNotPusher = errors.New("writer is not pusher")
)

func NewSender() *Sender {
	return &Sender{
		Chan: make(chan MessageInterface, 0),
	}
}

func (s *Sender) Send(useGzip bool, w http.ResponseWriter, request *http.Request) error {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var gzipWriter *gzip.Writer
	if useGzip {
		gzipWriter = gzip.NewWriter(w)
		w.Header().Set("Content-Encoding", "gzip")
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		return ErrWriterNotPusher
	}

	for {
		select {
		case <-request.Context().Done():
			fmt.Println("disconnected")
			return nil

		case msg := <-s.Chan:
			if useGzip {
				if _, err := io.WriteString(gzipWriter, EncodeMessage(msg)); err != nil {
					return fmt.Errorf("eventsource encode: %v", err)
				}

				gzipWriter.Flush()
			} else {
				if _, err := io.WriteString(w, EncodeMessage(msg)); err != nil {
					return fmt.Errorf("eventsource encode: %v", err)
				}
			}

			flusher.Flush()
		}
	}

	return nil
}
