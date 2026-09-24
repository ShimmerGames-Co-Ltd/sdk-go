package iap

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ShimmerGames-Co-Ltd/sdk-go/v3/core"
)

func TestVerifyOrderPathAndSuccess(t *testing.T) {
	t.Parallel()
	var gotPath, gotBody string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"code":0,"data":{"order_id":"o1","product_id":"sku","state":6}}`))
	}))
	t.Cleanup(srv.Close)
	corec, err := core.NewClient(core.WithURL(srv.URL), core.WithAppID("app1"), core.WithServerSignSecret("s"))
	if err != nil {
		t.Fatal(err)
	}
	cli, err := New(corec)
	if err != nil {
		t.Fatal(err)
	}
	rep, err := cli.VerifyOrder(context.Background(), VerifyOrderRequest{OrderID: "o1", Extras: "x"})
	if err != nil {
		t.Fatal(err)
	}
	if gotPath != pathVerifyOrder {
		t.Fatalf("path=%s", gotPath)
	}
	if !strings.Contains(gotBody, `"order_id"`) {
		t.Fatalf("body=%s", gotBody)
	}
	if !VerifySuccess(rep.State) || rep.ProductID != "sku" {
		t.Fatalf("%+v", rep)
	}
}

func TestRecordPurchaseAndHTTPError(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != pathRecordPurchase {
			t.Errorf("path=%s", r.URL.Path)
		}
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`no`))
	}))
	t.Cleanup(srv.Close)
	corec, _ := core.NewClient(core.WithURL(srv.URL), core.WithAppID("app1"), core.WithServerSignSecret("s"))
	cli, _ := New(corec)
	err := cli.RecordPurchase(context.Background(), RecordPurchaseRequest{PackageName: "pkg", PayChannel: "google", ProductID: "p", PurchaseToken: "t"})
	if _, ok := core.AsHTTPError(err); !ok {
		t.Fatalf("want HTTPError got %v", err)
	}
}
