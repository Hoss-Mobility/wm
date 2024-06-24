package internal

type SecretItem struct {
	Name               string `wm:"staff:r;developer:rw;admin:rw"`
	Comment            string `wm:"staff:rw;developer:rw;admin:rw"`
	SecretInfo         string `wm:"developer:r;admin:rw"`
	TopSecret          string `wm:"admin:rw"`
	CanOnlyBeWrittenTo string `wm:"staff:w;developer:w;admin:rw"`
}
