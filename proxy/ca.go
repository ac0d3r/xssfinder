package proxy

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	_ "embed"
	"io"
	"net/http"

	"github.com/elazarl/goproxy"
)

var (
	//go:embed xssfinder.ca.cert
	caCert []byte
	//go:embed xssfinder.ca.key
	caKey []byte
)

func SetCA(caCert, caKey []byte) error {
	goproxyCa, err := tls.X509KeyPair(caCert, caKey)
	if err != nil {
		return err
	}

	if goproxyCa.Leaf, err = x509.ParseCertificate(goproxyCa.Certificate[0]); err != nil {
		return err
	}

	goproxy.GoproxyCa = goproxyCa
	goproxy.OkConnect = &goproxy.ConnectAction{
		Action:    goproxy.ConnectAccept,
		TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.MitmConnect = &goproxy.ConnectAction{
		Action:    goproxy.ConnectMitm,
		TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.HTTPMitmConnect = &goproxy.ConnectAction{
		Action:    goproxy.ConnectHTTPMitm,
		TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.RejectConnect = &goproxy.ConnectAction{
		Action:    goproxy.ConnectReject,
		TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	return nil
}

func DownloadCaHandlerFunc(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	ctx.Logf("onRequest DownloadCaHanldeerFunc %s", req.URL.String())

	resp := &http.Response{
		StatusCode:    http.StatusOK,
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        make(http.Header),
		ContentLength: int64(len(caCert)),
		Body:          io.NopCloser(bytes.NewBuffer(caCert)),
	}
	resp.Header.Set("Content-Type", "application/x-x509-ca-cert")
	resp.Header.Set("Connection", "close")
	resp.Header.Set("Content-disposition", "attachment;filename=xxsfinder.cert")
	return req, resp
}
