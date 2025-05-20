module support

go 1.24.2

replace elog => ../elog

require (
	elog v0.0.0-00010101000000-000000000000
	github.com/minio/selfupdate v0.6.0
)

require (
	aead.dev/minisign v0.3.0 // indirect
	golang.org/x/crypto v0.38.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
)
