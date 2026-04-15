package meta_file_persistence

import "github.com/janicaleksander/cloud/claimservice/domain"

func MetaFileModelToDomain(f *MetaFileModel) *domain.MetaFile {
	return &domain.MetaFile{
		ID:       f.ID,
		FileName: f.FileName,
		FileExt:  f.FileExt,
		FileSize: f.FileSize,
		Date:     f.Date,
		FileURL:  f.FileURL,
	}
}

func MetaFileDomainToModel(f *domain.MetaFile) *MetaFileModel {
	return &MetaFileModel{
		ID:       f.ID,
		FileName: f.FileName,
		FileExt:  f.FileExt,
		FileSize: f.FileSize,
		Date:     f.Date,
		FileURL:  f.FileURL,
	}
}
