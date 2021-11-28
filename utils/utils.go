package utils

import (
	"errors"
	"google.golang.org/grpc"
	"io"
	"log"
)

func CreateBytesResponse(clientStream grpc.ClientStream, reqBytes []byte) ([]byte, error) {
	errCh1 := SendRequestIntoStream(clientStream, reqBytes)
	errCh2, retChan := RetrieveResponseFromStream(clientStream)
	for i := 0; i < 3; i++ {
		select {
		case magErr := <-errCh1:
			if !errors.Is(magErr, io.EOF) {
				return nil, magErr
			}
			clientStream.CloseSend()
		case magErr2 := <-errCh2:
			if !errors.Is(magErr2, io.EOF) {
				return nil, magErr2
			}
		case response := <-retChan:
			return response, nil
		}
	}
	log.Println("SOMETHING DEFAULT HAPPENED:")

	return nil, nil
}

func SendRequestIntoStream(stream grpc.ClientStream, info []byte) chan error {
	ret := make(chan error, 1)
	go func() {
		for {
			if err := stream.SendMsg(&info); err != nil {
				ret <- err

				break
			}
		}

	}()
	return ret
}

func RetrieveResponseFromStream(stream grpc.ClientStream) (chan error, chan []byte) {
	errCh := make(chan error, 1)
	retChan := make(chan []byte, 1)
	var a []byte
	go func() {
		for {
			if err := stream.RecvMsg(&a); err != nil {
				if errors.Is(err, io.EOF) {
					retChan <- a
				}
				errCh <- err
				break
			}
		}

	}()

	return errCh, retChan
}
