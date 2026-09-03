package middleware

import (
	"testing"
)

// ============ ParseUserAgent ============

func TestParseUserAgentBot(t *testing.T) {
	device, browser, os := ParseUserAgent("Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")
	if device != "bot" || browser != "bot" || os != "bot" {
		t.Errorf("expected bot/bot/bot, got %s/%s/%s", device, browser, os)
	}
}

func TestParseUserAgentChromeWindows(t *testing.T) {
	device, browser, os := ParseUserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	if device != "desktop" {
		t.Errorf("expected desktop, got %s", device)
	}
	if browser != "Chrome" {
		t.Errorf("expected Chrome, got %s", browser)
	}
	if os != "Windows" {
		t.Errorf("expected Windows, got %s", os)
	}
}

func TestParseUserAgentMobileAndroid(t *testing.T) {
	device, browser, os := ParseUserAgent("Mozilla/5.0 (Linux; Android 14; Pixel 8) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36")
	if device != "mobile" {
		t.Errorf("expected mobile, got %s", device)
	}
	if browser != "Chrome" {
		t.Errorf("expected Chrome, got %s", browser)
	}
	if os != "Android" {
		t.Errorf("expected Android, got %s", os)
	}
}

func TestParseUserAgentSafariMacOS(t *testing.T) {
	device, browser, os := ParseUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.2 Safari/605.1.15")
	if device != "desktop" {
		t.Errorf("expected desktop, got %s", device)
	}
	if browser != "Safari" {
		t.Errorf("expected Safari, got %s", browser)
	}
	if os != "macOS" {
		t.Errorf("expected macOS, got %s", os)
	}
}

func TestParseUserAgentIPad(t *testing.T) {
	device, _, os := ParseUserAgent("Mozilla/5.0 (iPad; CPU OS 17_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.2 Mobile/15E148 Safari/604.1")
	if device != "tablet" {
		t.Errorf("expected tablet, got %s", device)
	}
	if os != "iOS" {
		t.Errorf("expected iOS, got %s", os)
	}
}

func TestParseUserAgentFirefoxLinux(t *testing.T) {
	device, browser, os := ParseUserAgent("Mozilla/5.0 (X11; Linux x86_64; rv:120.0) Gecko/20100101 Firefox/120.0")
	if device != "desktop" {
		t.Errorf("expected desktop, got %s", device)
	}
	if browser != "Firefox" {
		t.Errorf("expected Firefox, got %s", browser)
	}
	if os != "Linux" {
		t.Errorf("expected Linux, got %s", os)
	}
}

func TestParseUserAgentEdge(t *testing.T) {
	_, browser, _ := ParseUserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0")
	if browser != "Edge" {
		t.Errorf("expected Edge, got %s", browser)
	}
}

func TestParseUserAgentEmpty(t *testing.T) {
	device, browser, os := ParseUserAgent("")
	if device != "unknown" || browser != "unknown" || os != "unknown" {
		t.Errorf("expected unknown/unknown/unknown, got %s/%s/%s", device, browser, os)
	}
}
// ============ classifyVisit ============

func TestClassifyVisit(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"/api/v1/public/products/custom-leggings", "page_view"},
		{"/api/v1/public/blogs/hello", "page_view"},
		{"/api/v1/public/pages/about", "page_view"},
		{"/api/v1/public/cases/test", "page_view"},
		{"/api/v1/public/home", "page_view"},
		{"/api/v1/public/about", "page_view"},
		{"/api/v1/public/faq", "page_view"},
		{"/api/v1/public/contact", "page_view"},
		{"/api/v1/public/categories", "page_view"},
		{"/api/v1/public/series", "page_view"},
		{"/api/v1/public/fabrics", "page_view"},
		{"/api/v1/public/factories", "page_view"},
		{"/api/v1/public/certifications", "page_view"},
		{"/api/v1/public/navigations", "api_call"},
		{"/api/v1/public/theme", "api_call"},
		{"/api/v1/public/currencies", "api_call"},
		{"/api/v1/public/geo", "api_call"},
		{"/api/v1/public/i18n", "api_call"},
		{"/api/v1/public/locale", "api_call"},
		{"/health", "api_call"},
	}
	for _, tt := range tests {
		got := classifyVisit(tt.path)
		if got != tt.want {
			t.Errorf("classifyVisit(%q) = %q, want %q", tt.path, got, tt.want)
		}
	}
}

// ============ parseSource ============

func TestParseSource(t *testing.T) {
	tests := []struct {
		referer     string
		utmSource   string
		utmMedium   string
		utmCampaign string
		wantSource  string
		wantMedium  string
	}{
		// direct (no referer, no utm)
		{"", "", "", "", "direct", "none"},
		// organic search
		{"https://www.google.com/search?q=sportswear", "", "", "", "google", "organic"},
		{"https://www.bing.com/search?q=yoga+pants", "", "", "", "bing", "organic"},
		{"https://www.baidu.com/link?url=xxx", "", "", "", "baidu", "organic"},
		// social
		{"https://www.facebook.com/some-post", "", "", "", "facebook", "social"},
		{"https://www.linkedin.com/feed/", "", "", "", "linkedin", "social"},
		{"https://twitter.com/some-tweet", "", "", "", "twitter", "social"},
		{"https://x.com/some-post", "", "", "", "twitter", "social"},
		{"https://www.youtube.com/watch?v=abc", "", "", "", "youtube", "social"},
		{"https://www.instagram.com/p/xyz", "", "", "", "instagram", "social"},
		// referral (other referer)
		{"https://some-blog.com/article", "", "", "", "referral", "referral"},
		// UTM overrides
		{"", "newsletter", "email", "spring_sale", "newsletter", "email"},
		{"https://www.google.com/", "google", "cpc", "campaign_a", "google", "cpc"},
		// UTM medium overrides parsed medium
		{"https://www.facebook.com/post", "facebook", "paid_social", "summer", "facebook", "paid_social"},
	}
	for _, tt := range tests {
		source, medium, _ := parseSource(tt.referer, tt.utmSource, tt.utmMedium, tt.utmCampaign)
		if source != tt.wantSource || medium != tt.wantMedium {
			t.Errorf("parseSource(%q, %q, %q, %q) = (%q, %q, _), want (%q, %q, _)",
				tt.referer, tt.utmSource, tt.utmMedium, tt.utmCampaign, source, medium, tt.wantSource, tt.wantMedium)
		}
	}
}

// ============ AnonymizeIP ============

func TestAnonymizeIP(t *testing.T) {
	tests := []struct {
		ip   string
		want string
	}{
		{"192.168.1.42", "192.168.1.0"},
		{"10.0.0.5", "10.0.0.0"},
		{"203.0.113.99", "203.0.113.0"},
		{"", "anonymous"},
		{"unknown", "anonymous"},
		// IPv6 simplified test：仅保留前 3 段 + ::/0 后缀（当前实现所产生的结果）
		{"2001:db8::1", "2001:db8:::/0"},
		{"2001:db8:abcd:1234:5678::1", "2001:db8:abcd::/0"},
	}
	for _, tt := range tests {
		got := AnonymizeIP(tt.ip)
		if got != tt.want {
			t.Errorf("AnonymizeIP(%q) = %q, want %q", tt.ip, got, tt.want)
		}
	}
}

// ============ firstNonEmpty ============

func TestFirstNonEmpty(t *testing.T) {
	if got := firstNonEmpty("", "b", "c"); got != "b" {
		t.Errorf("expected 'b', got %q", got)
	}
	if got := firstNonEmpty("a", "b"); got != "a" {
		t.Errorf("expected 'a', got %q", got)
	}
	if got := firstNonEmpty(""); got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}