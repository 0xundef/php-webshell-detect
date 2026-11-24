package lang

import "testing"

func TestMethodLookupIsCaseInsensitive(t *testing.T) {
	class := CreatePHPClass("EntryMain", "main")
	method := CreatePHPMethod(class, "Shell_Exec", nil, nil, []string{"$sink"})
	class.AddMethod(method)

	if got := class.GetDeclaredMethod("shell_exec"); got != method {
		t.Fatalf("GetDeclaredMethod returned %v, want the declared method", got)
	}

	h := CreateClassHierarchy()
	h.AddClass(class)
	if got := h.LookupMethod(class, "SHELL_EXEC"); got != method {
		t.Fatalf("LookupMethod returned %v, want the declared method", got)
	}
}

func TestClassLookupIsCaseInsensitiveButFieldsAreNot(t *testing.T) {
	class := CreatePHPClass("MyClass", "ns")
	h := CreateClassHierarchy()
	h.AddClass(class)
	if got := h.GetClass("myclass"); got != class {
		t.Fatalf("GetClass returned %v, want %v", got, class)
	}

	declared := &PHPClassField{DeclaringClass: class, Name: "Cmd"}
	class.AddField(declared)
	if got := class.GetDeclaredField("Cmd"); got != declared {
		t.Fatalf("GetDeclaredField(Cmd) = %v, want the declared field", got)
	}
	if got := class.GetDeclaredField("cmd"); got == declared {
		t.Fatal("field lookup must stay case-sensitive")
	}
}
