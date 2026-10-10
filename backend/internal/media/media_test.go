package media

import (
	"strings"
	"testing"
)

func TestSanitize(t *testing.T) {
	cases := map[string]string{
		"Profile Photo.JPG":      "profile-photo.jpg",
		"../../etc/passwd.png":   "passwd.png",
		"C:\\Users\\me\\cv.pdf":  "cv.pdf",
		"weird  name!!.jpeg":     "weird-name.jpg",
		"...png":                 "file.png",
		"résumé.pdf":             "r-sum.pdf",
		"a.b.c.webp":             "a.b.c.webp",
	}
	for in, want := range cases {
		if got := sanitize(in); got != want {
			t.Errorf("sanitize(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestSaveListDelete(t *testing.T) {
	lib, err := New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := lib.Save("evil.exe", strings.NewReader("x")); err == nil {
		t.Fatal("expected disallowed extension error")
	}
	f1, err := lib.Save("Photo.PNG", strings.NewReader("one"))
	if err != nil {
		t.Fatal(err)
	}
	f2, err := lib.Save("photo.png", strings.NewReader("two"))
	if err != nil {
		t.Fatal(err)
	}
	if f1.Name != "photo.png" || f2.Name != "photo-2.png" {
		t.Fatalf("unexpected names %q %q", f1.Name, f2.Name)
	}
	if !f1.IsImage || f1.URL != "/uploads/photo.png" {
		t.Fatalf("unexpected describe: %+v", f1)
	}
	files, err := lib.List()
	if err != nil || len(files) != 2 {
		t.Fatalf("list: %v %d", err, len(files))
	}
	if err := lib.Delete("../photo.png"); err == nil {
		t.Fatal("expected invalid name error")
	}
	if err := lib.Delete("photo.png"); err != nil {
		t.Fatal(err)
	}
	if err := lib.Delete("photo.png"); err == nil {
		t.Fatal("expected not found on second delete")
	}
}
