package services

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type Ad struct {
    Date          string                 `json:"date" firestore:"Date"`
    Description   string                 `json:"description" firestore:"Description"`
    ParentId      *firestore.DocumentRef `json:"parentId" firestore:"ParentId"` // Keep as DocumentRef for Firestore
    Salary        int                    `json:"salary" firestore:"Salary"`
    School        string                 `json:"school" firestore:"School"`
    SchoolAddress string                 `json:"schoolAddress" firestore:"SchoolAdress"`
    Status        string                 `json:"status" firestore:"Status"`
    Time          time.Time              `json:"time" firestore:"Time"`
}

// AdResponse is used to format the response for JSON output
type AdResponse struct {
    ID            string `json:"id"`            // Document ID for each ad
    Date          string `json:"date"`
    Description   string `json:"description"`
    ParentId      string `json:"parentId"`
    Salary        int    `json:"salary"`
    School        string `json:"school"`
    SchoolAddress string `json:"schoolAddress"`
    Status        string `json:"status"`
    Time          string `json:"time"`
}


type AdService struct {
    FirestoreClient *firestore.Client
}

func NewAdService(client *firestore.Client) *AdService {
    return &AdService{FirestoreClient: client}
}

// FetchAllAds retrieves all ads from Firestore
func (s *AdService) FetchAllAds(ctx context.Context) ([]AdResponse, error) {
    var adResponses []AdResponse
    iter := s.FirestoreClient.Collection("ad").Documents(ctx)

    log.Println("Starting to fetch documents from the 'ad' collection...")

    for {
        doc, err := iter.Next()
        if err == iterator.Done {
            log.Println("Finished fetching all documents.")
            break
        }
        if err != nil {
            log.Printf("Error fetching ad document: %v", err)
            return nil, err
        }

        log.Printf("Document found with ID: %s", doc.Ref.ID)

        // Map Firestore document to Ad struct
        var ad Ad
        if err := doc.DataTo(&ad); err != nil {
            log.Printf("Error converting document ID %s to Ad struct: %v", doc.Ref.ID, err)
            return nil, err
        }

        // Convert to AdResponse format
        adResponse := AdResponse{
            ID:            doc.Ref.ID,                  // Set document ID here
            Date:          ad.Date,
            Description:   ad.Description,
            Salary:        ad.Salary,
            School:        ad.School,
            SchoolAddress: ad.SchoolAddress,
            Status:        ad.Status,
            Time:          ad.Time.Format(time.RFC3339),
        }

        // Convert ParentId (DocumentRef) to string if it's not nil
        if ad.ParentId != nil {
            adResponse.ParentId = ad.ParentId.Path // Convert DocumentRef to its path as a string
        } else {
            adResponse.ParentId = ""
        }

        log.Printf("Ad data retrieved and formatted: %+v", adResponse)

        adResponses = append(adResponses, adResponse)
    }
    return adResponses, nil
}

