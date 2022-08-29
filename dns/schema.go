package dns

import "github.com/jobstoit/hcloud-dns-go/dns/schema"

// ZoneFromSchema converts schema.Zone to a Zone.
func ZoneFromSchema(s schema.Zone) *Zone {
	zone := &Zone{
		ID:              s.ID,
		Created:         s.Created,
		Modified:        s.Modified,
		LegacyDNSHost:   s.LegacyDNSHost,
		LegacyNS:        s.LegacyNS,
		Name:            s.Name,
		NS:              s.NS,
		Owner:           s.Owner,
		Paused:          s.Paused,
		Permission:      s.Permission,
		Project:         s.Project,
		Registrar:       s.Registrar,
		Status:          ZoneStatus(s.Status),
		Ttl:             s.Ttl,
		Verified:        s.Verified,
		RecordsCount:    s.RecordsCount,
		IsSecondaryDNS:  s.IsSecondaryDNS,
		TxtVerification: TxtVerificationFromSchema(s.TxtVerification),
	}

	return zone
}

// TxtVerificationFromSchema converts schema.TxtVerification to TxtVerification
func TxtVerificationFromSchema(s schema.TxtVerification) *TxtVerification {
	return &TxtVerification{
		Name:  s.Name,
		Token: s.Token,
	}
}

// PaginationFromSchema converts a schema.MetaPagination to a Pagination.
func PaginationFromSchema(s schema.MetaPagination) Pagination {
	return Pagination{
		LastPage:     s.LastPage,
		Page:         s.Page,
		PerPage:      s.PerPage,
		TotalEntries: s.TotalEntries,
	}
}

// RecordFromSchema convers a schema.Record to Record
func RecordFromSchema(s schema.Record) *Record {
	return &Record{
		Type:     RecordType(s.Type),
		ID:       s.ID,
		Created:  s.Created,
		Modified: s.Modified,
		Zone:     &Zone{ID: s.ZoneID},
		Name:     s.Name,
		Value:    s.Value,
		Ttl:      s.Ttl,
	}
}

// ValidateZoneFileFromSchema converts a schema.ValidateZoneFileResponse to a ValidatedZoneFile.
func ValidateZoneFileFromSchema(s schema.ValidateZoneFileResponse) *ValidatedZoneFile {
	zf := &ValidatedZoneFile{
		PassedRecords: s.PassedRecords,
		ValidRecords:  []*Record{},
	}

	for _, rec := range s.ValidRecords {
		zf.ValidRecords = append(zf.ValidRecords, RecordFromSchema(rec))
	}

	return zf
}

// RecordEntryFromSchema converts a schema.RecordCreateRequest to a RecordCreateOpts
func RecordEntryFromSchema(s schema.RecordBulkEntry) *RecordEntry {
	return &RecordEntry{
		Type:   RecordType(s.Type),
		ZoneID: s.ZoneID,
		Name:   s.Name,
		Value:  s.Value,
		Ttl:    s.Ttl,
	}
}

// PrimaryServerFromSchema converts a schema.PrimaryServer to a PrimaryServer
func PrimaryServerFromSchema(s schema.PrimaryServer) *PrimaryServer {
	return &PrimaryServer{
		ID:       s.ID,
		Port:     s.Port,
		Created:  s.Created,
		Modified: s.Modified,
		Zone:     &Zone{ID: s.ZoneID},
		Address:  s.Address,
	}
}
