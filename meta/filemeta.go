package meta

// FileMeta: struct for file meta info
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}

// UpdateFileMeta: Add/update file meta
func UpdateFileMeta(fmeta FileMeta) {
	fileMetas[fmeta.FileSha1] = fmeta
}

// GetFileMeta: get the file meta from sha1
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}
