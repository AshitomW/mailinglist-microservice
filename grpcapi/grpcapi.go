package grpcapi

import (
	"database/sql"
	"time"

	"mailinglist-ms/maildb"
	pb "mailinglist-ms/proto"
)

type MailServer struct {
	pb.UnimplementedMailingListServiceServer
	db *sql.DB
}

// conver protbuf message

func pbEntryToMailDbEntry(pbEntry *pb.EmailEntry) maildb.EmailEntry{
	tme := time.Unix(pbEntry.ConfirmedAt,0)
	return maildb.EmailEntry{
		Id: pbEntry.Id,
		Email: pbEntry.Email,
		ConfirmedAt: &tme,
		OptOut: pbEntry.OptOut,
	}

}


func mailDbToPbEntry(maildbEntry *maildb.EmailEntry)pb.EmailEntry{
	return pb.EmailEntry{
		Id: maildbEntry.Id,
		Email: maildbEntry.Email,
		ConfirmedAt: maildbEntry.ConfirmedAt.Unix(),
		OptOut: maildbEntry.OptOut,
	}
}
