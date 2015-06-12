package flash

import "testing"

func TestMakeActions(t *testing.T) {
	//  Standard CRUD
	assertEqual(t, "Index", makeAction("GET", "", "", []string{}))
	assertEqual(t, "Create", makeAction("POST", "", "", []string{}))
	assertEqual(t, "Show", makeAction("GET", "1", "", []string{}))
	assertEqual(t, "Update", makeAction("POST", "1", "", []string{}))
	assertEqual(t, "Update", makeAction("PUT", "1", "", []string{}))
	assertEqual(t, "Destroy", makeAction("DELETE", "1", "", []string{}))
	// ID with actions
	assertEqual(t, "GETAction", makeAction("GET", "1", "action", []string{}))
	assertEqual(t, "POSTAction", makeAction("POST", "1", "action", []string{}))
	assertEqual(t, "PUTAction", makeAction("PUT", "1", "action", []string{}))
	assertEqual(t, "DELETEAction", makeAction("DELETE", "1", "action", []string{}))
	// Extra functions
	assertEqual(t, "GETAction", makeAction("GET", "action", "", []string{"GETAction"}))
	assertEqual(t, "POSTAction", makeAction("POST", "action", "", []string{"POSTAction"}))
}

func BenchmarkMakeActionIndex(b *testing.B) {
	for n := 0; n < b.N; n++ {
		makeAction("GET", "", "", []string{})
	}
}

func BenchmarkMakeActionCreate(b *testing.B) {
	for n := 0; n < b.N; n++ {
		makeAction("POST", "", "", []string{})
	}
}

func BenchmarkMakeActionShowww(b *testing.B) {
	for n := 0; n < b.N; n++ {
		makeAction("GET", "1", "", []string{})
	}
}

func BenchmarkMakeActionUpdate(b *testing.B) {
	for n := 0; n < b.N; n++ {
		makeAction("PUT", "1", "", []string{})
	}
}

func BenchmarkMakeActionDelete(b *testing.B) {
	for n := 0; n < b.N; n++ {
		makeAction("DELETE", "1", "", []string{})
	}
}

func BenchmarkMakeActionGETact(b *testing.B) {
	for n := 0; n < b.N; n++ {
		makeAction("GET", "1", "action", []string{})
	}
}
func BenchmarkMakeActionGETac2(b *testing.B) {
	for n := 0; n < b.N; n++ {
		makeAction("GET", "action", "", []string{"GETAction", "GETAction1", "GETAction2", "GETAction3"})
	}
}
func BenchmarkMakeActionGETac3(b *testing.B) {
	for n := 0; n < b.N; n++ {
		makeAction("GET", "action3", "", []string{"GETAction", "GETAction1", "GETAction2", "GETAction3"})
	}
}
