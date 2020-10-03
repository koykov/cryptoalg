package doubleclick

import (
	"bytes"
	"testing"
)

var (
	encryptedAdID = []byte{
		0x38, 0x6e, 0x3a, 0xc0, 0x00, 0x0c, 0x0a, 0x08, 0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0xe9, 0x8c,
		0xc9, 0x46, 0x57, 0x3f, 0xbf, 0x46, 0x45, 0xef, 0x06, 0x0b, 0x17, 0xa6, 0x67, 0xa6, 0x17, 0xc6, 0x6b, 0xcb,
	}
	decryptedAdID   = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
	decryptedAdUUID = []byte("00010203-0405-0607-0809-0a0b0c0d0e0f")
)

func TestAdID_Decrypt(t *testing.T) {
	d := NewAdID(encryptionKey, integrityKey)
	var (
		dst []byte
		err error
	)
	dst, err = d.Decrypt(dst, encryptedAdID, false)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(dst, decryptedAdID) {
		t.Error("decrypt AdID failed")
	}
	dst, _ = d.Decrypt(dst[:0], encryptedAdID, true)
	if !bytes.Equal(dst, decryptedAdUUID) {
		t.Error("decrypt AdID UUID failed")
	}
}

func BenchmarkAdID_Decrypt(b *testing.B) {
	d := NewAdID(encryptionKey, integrityKey)
	var (
		dst []byte
		err error
	)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		dst = dst[:0]
		dst, err = d.Decrypt(dst, encryptedAdID, true)
		if err != nil {
			b.Error(err)
		}
		if !bytes.Equal(dst, decryptedAdUUID) {
			b.Error("decrypt AdID failed")
		}
		d.Reset()
	}
}

func benchmarkAdIDDecryptParallel(b *testing.B, n int) {
	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		var (
			dst []byte
			err error
		)
		for pb.Next() {
			for i := 0; i < n; i++ {
				a := AcquireAdID(encryptionKey, integrityKey)

				dst = dst[:0]
				dst, err = a.Decrypt(dst, encryptedAdID, true)
				if err != nil {
					b.Error(err)
				}
				if !bytes.Equal(dst, decryptedAdUUID) {
					b.Error("decrypt AdID failed")
				}

				ReleaseAdID(a)
			}
		}
	})
}

func BenchmarkAdID_DecryptParallel1(b *testing.B) {
	benchmarkAdIDDecryptParallel(b, 1)
}

func BenchmarkAdID_DecryptParallel10(b *testing.B) {
	benchmarkAdIDDecryptParallel(b, 10)
}

func BenchmarkAdID_DecryptParallel100(b *testing.B) {
	benchmarkAdIDDecryptParallel(b, 100)
}

func BenchmarkAdID_DecryptParallel1000(b *testing.B) {
	benchmarkAdIDDecryptParallel(b, 1000)
}
