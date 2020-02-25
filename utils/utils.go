package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/lambda"
)

func CheckError(err error) {
	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		case lambda.ErrCodeServiceException:
			fmt.Println(lambda.ErrCodeServiceException, aerr.Error())
		case lambda.ErrCodeTooManyRequestsException:
			fmt.Println(lambda.ErrCodeTooManyRequestsException, aerr.Error())
		case lambda.ErrCodeInvalidParameterValueException:
			fmt.Println(lambda.ErrCodeInvalidParameterValueException, aerr.Error())
		default:
			fmt.Println(aerr.Error())
		}
	} else {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
	}
}

//from https://stackoverflow.com/questions/20357223/easy-way-to-unzip-file-with-golang
func Unzip(zipFile, destination string) error {
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()
	os.MkdirAll(destination, 0755)
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()
		path := filepath.Join(destination, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()
			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}
	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return nil
		}
	}
	return nil
}
