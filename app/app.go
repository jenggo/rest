package app

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"

	"rest/app/exchange/build"
	"rest/app/flags"
	"rest/app/input"
	"rest/vars"

	"rest/app/exchange"

	"rest/app/output"

	"github.com/pkg/errors"
)

func Run(options *vars.Options) error {
	// Parse flags
	args, usage, optionSet, err := flags.Parse(os.Args)
	if err != nil {
		return err
	}
	inputOptions := optionSet.InputOptions
	exchangeOptions := optionSet.ExchangeOptions
	exchangeOptions.Transport = options.Transport
	outputOptions := optionSet.OutputOptions

	// Parse positional arguments
	in, err := input.ParseArgs(args, os.Stdin, &inputOptions)
	if _, ok := errors.Cause(err).(*input.UsageError); ok {
		usage.PrintUsage(os.Stderr)
		return err
	}
	if err != nil {
		return err
	}

	// Send request and receive response
	status, err := Exchange(in, &exchangeOptions, &outputOptions)
	if err != nil {
		return err
	}

	if exchangeOptions.CheckStatus {
		os.Exit(getExitStatus(status))
	}

	return nil
}

func getExitStatus(statusCode int) int {
	if 300 <= statusCode && statusCode < 600 {
		return statusCode / 100
	}
	return 0
}

func Exchange(in *input.Input, exchangeOptions *vars.ExchangeOptions, outputOptions *vars.OutputOptions) (int, error) {
	// Prepare printer
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	printer := output.NewPrinter(writer, outputOptions)

	// Build HTTP request
	request, err := build.HTTPRequest(in, exchangeOptions)
	if err != nil {
		return -1, err
	}

	// Print HTTP request
	if outputOptions.PrintRequestHeader || outputOptions.PrintRequestBody {
		// `request` does not contain HTTP headers that HttpClient.Do adds.
		// We can get these headers by DumpRequestOut and ReadRequest.
		dump, err := httputil.DumpRequestOut(request, true)
		if err != nil {
			return -1, err // should not happen
		}
		r, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(dump)))
		if err != nil {
			return -1, err // should not happen
		}
		defer r.Body.Close()

		// ReadRequest deletes Host header. We must restore it.
		if request.Host != "" {
			r.Header.Set("Host", request.Host)
		} else {
			r.Header.Set("Host", request.URL.Host)
		}

		if outputOptions.PrintRequestHeader {
			if err := printer.PrintRequestLine(r); err != nil {
				return -1, err
			}
			if err := printer.PrintHeader(r.Header); err != nil {
				return -1, err
			}
		}
		if outputOptions.PrintRequestBody {
			if err := printer.PrintBody(r.Body, r.Header.Get("Content-Type")); err != nil {
				return -1, err
			}
		}
		fmt.Fprintln(writer)
		writer.Flush()
	}

	// Send HTTP request and receive HTTP request
	httpClient, err := exchange.BuildHTTPClient(exchangeOptions)
	if err != nil {
		return -1, err
	}
	resp, err := httpClient.Do(request)
	if err != nil {
		return -1, errors.Wrap(err, "sending HTTP request")
	}
	defer resp.Body.Close()

	if outputOptions.PrintResponseHeader {
		if err := printer.PrintStatusLine(resp.Proto, resp.Status, resp.StatusCode); err != nil {
			return -1, err
		}
		if err := printer.PrintHeader(resp.Header); err != nil {
			return -1, err
		}
		writer.Flush()
	}

	if outputOptions.Download {
		file := output.NewFileWriter(in.URL, outputOptions)

		if err := printer.PrintDownload(resp.ContentLength, file.Filename()); err != nil {
			return -1, err
		}
		writer.Flush()

		if err = file.Download(resp); err != nil {
			return -1, err
		}
	} else {
		if outputOptions.PrintResponseBody {
			if err := printer.PrintBody(resp.Body, resp.Header.Get("Content-Type")); err != nil {
				return -1, err
			}
		}
	}

	return resp.StatusCode, nil
}
