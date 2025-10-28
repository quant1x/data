package exchange

import "testing"

func TestGetMarket_Cases(t *testing.T) {
    cases := []struct{
        symbol string
        want   string
    }{
        {"600000", "sh"},
        {"600000.SH", "sh"},
        {"SZ000002", "sz"},
        {"000001", "sh"},
        {"002001", "sz"},
        {"300001", "sz"},
        {"399001", "sz"},
        {"899001", "bj"},
        {"012345.HK", "hk"},
        {"USAAPL", "us"},
        {"AAPL.US", "us"},
    }

    for _, c := range cases {
        got := GetMarket(c.symbol)
        if got != c.want {
            t.Errorf("GetMarket(%q) = %q; want %q", c.symbol, got, c.want)
        }
    }
}

func TestDetectMarketAndCorrectSecurityCode(t *testing.T) {
    // DetectMarket with suffix form
    mid, flag, code := DetectMarket("600000.SH")
    if flag != "sh" || mid != MarketIdShangHai || code != "600000" {
        t.Fatalf("DetectMarket(600000.SH) = (%v,%q,%q); want (MarketIdShangHai, 'sh', '600000')", mid, flag, code)
    }

    // CorrectSecurityCode should normalize to flag+code
    if got := CorrectSecurityCode("600000.SH"); got != "sh600000" {
        t.Fatalf("CorrectSecurityCode(600000.SH) = %q; want %q", got, "sh600000")
    }

    // Beijing index
    mid2, flag2, code2 := DetectMarket("899001")
    if flag2 != "bj" || mid2 != MarketIdBeiJing || code2[:3] != "899" {
        t.Fatalf("DetectMarket(899001) = (%v,%q,%q)", mid2, flag2, code2)
    }

    // GetMarketId mapping
    if id := GetMarketId("899001"); id != MarketIdBeiJing {
        t.Fatalf("GetMarketId(899001) = %v; want %v", id, MarketIdBeiJing)
    }
}

func TestAssertIndexAndStockBySecurityCode(t *testing.T) {
    if !AssertIndexBySecurityCode("000001") {
        t.Fatalf("AssertIndexBySecurityCode(000001) = false; want true")
    }

    if !AssertIndexBySecurityCode("899001") {
        t.Fatalf("AssertIndexBySecurityCode(899001) = false; want true")
    }

    if !AssertStockBySecurityCode("600000.SH") {
        t.Fatalf("AssertStockBySecurityCode(600000.SH) = false; want true")
    }

    // Ensure indexes are not classified as stocks
    if AssertStockBySecurityCode("000001.SH") {
        t.Fatalf("AssertStockBySecurityCode(000001.SH) = true; want false")
    }
}
