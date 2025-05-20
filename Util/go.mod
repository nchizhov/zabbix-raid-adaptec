module util

go 1.24.2

replace support => ./support

replace adapter => ./adapter

replace elog => ./elog

require (
	adapter v0.0.0-00010101000000-000000000000 // indirect
	aead.dev/minisign v0.2.0 // indirect
	elog v0.0.0-00010101000000-000000000000 // indirect
	github.com/minio/selfupdate v0.6.0 // indirect
	golang.org/x/crypto v0.35.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	support v0.0.0-00010101000000-000000000000 // indirect
)
