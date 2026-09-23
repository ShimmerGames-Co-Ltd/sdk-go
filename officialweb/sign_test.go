package officialweb

import "testing"

func TestSignSortedQSMD5LookupVector(t *testing.T) {
	t.Parallel()
	got := SignSortedQSMD5(map[string]any{"project_role_id": "player-web-prod"}, "test-sign-secret")
	const want = "de2d26ef096be852ede87dc205eb59e9"
	if got != want {
		t.Fatalf("got %s want %s", got, want)
	}
}

func TestSignSortedQSMD5EmptyBody(t *testing.T) {
	t.Parallel()
	got := SignSortedQSMD5(map[string]any{}, "test-sign-secret")
	const want = "cfea7ff286f55e32220404f77a34f1e1"
	if got != want {
		t.Fatalf("got %s want %s", got, want)
	}
}

func TestSignSortedQSMD5StripsSignKey(t *testing.T) {
	t.Parallel()
	a := SignSortedQSMD5(map[string]any{"project_role_id": "player-web-prod", "sign": "ignored"}, "test-sign-secret")
	b := SignSortedQSMD5(map[string]any{"project_role_id": "player-web-prod"}, "test-sign-secret")
	if a != b {
		t.Fatalf("sign key should be ignored: %s vs %s", a, b)
	}
}
