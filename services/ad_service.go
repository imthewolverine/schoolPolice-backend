package services

import (
	"context"
	"fmt"
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

// FetchAdByID retrieves a single ad by document ID
func (s *AdService) FetchAdByID(ctx context.Context, id string) (*AdResponse, error) {
    docRef := s.FirestoreClient.Collection("ad").Doc(id)
    doc, err := docRef.Get(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to retrieve ad: %v", err)
    }

    var ad Ad
    if err := doc.DataTo(&ad); err != nil {
        return nil, fmt.Errorf("failed to convert document data to Ad struct: %v", err)
    }

    // Create and return the AdResponse struct with document ID and formatted data
    adResponse := &AdResponse{
        ID:            doc.Ref.ID,
        Date:          ad.Date,
        Description:   ad.Description,
        Salary:        ad.Salary,
        School:        ad.School,
        SchoolAddress: ad.SchoolAddress,
        Status:        ad.Status,
        Time:          ad.Time.Format(time.RFC3339),
    }

    if ad.ParentId != nil {
        adResponse.ParentId = ad.ParentId.Path // Convert DocumentRef to path
    }

    return adResponse, nil
}

func (s *AdService) AddAd(ctx context.Context, ad Ad) (string, error) {
    docRef, _, err := s.FirestoreClient.Collection("ad").Add(ctx, map[string]interface{}{
        "Date":          ad.Date,
        "Description":   ad.Description,
        "ParentId":      ad.ParentId,
        "Salary":        ad.Salary,
        "School":        ad.School,
        "SchoolAdress":  ad.SchoolAddress,
        "Status":        ad.Status,
        "Time":          ad.Time,
    })
    if err != nil {
        return "", err
    }
    return docRef.ID, nil
}


func (s *AdService) DeleteAdByID(ctx context.Context, id string) error {
    _, err := s.FirestoreClient.Collection("ad").Doc(id).Delete(ctx)
    return err
}
func (s *AdService) UpdateAdByID(ctx context.Context, id string, ad Ad) error {
    _, err := s.FirestoreClient.Collection("ad").Doc(id).Set(ctx, map[string]interface{}{
        "Date":          ad.Date,
        "Description":   ad.Description,
        "ParentId":      ad.ParentId,
        "Salary":        ad.Salary,
        "School":        ad.School,
        "SchoolAdress":  ad.SchoolAddress,
        "Status":        ad.Status,
        "Time":          ad.Time,
    }, firestore.MergeAll) // MergeAll updates only the provided fields
    return err
}

func (s *AdService) CreateRequest(ctx context.Context, workerID string) (*firestore.DocumentRef, error) {
    requestRef, _, err := s.FirestoreClient.Collection("requests").Add(ctx, map[string]interface{}{
        "Status":   "",                  // Set the initial status
        "WorkerID": s.FirestoreClient.Doc(workerID), // Convert workerID string to a DocumentRef
    })
    if err != nil {
        return nil, err
    }
    return requestRef, nil
}
func (s *AdService) AddRequestToAd(ctx context.Context, adID string, requestRef *firestore.DocumentRef) error {
    adRef := s.FirestoreClient.Collection("ad").Doc(adID)
    _, err := adRef.Update(ctx, []firestore.Update{
        {
            Path:  "requests",
            Value: firestore.ArrayUnion(requestRef), // Add to requests array
        },
    })
    return err
}

func (s *AdService) GetUserAdsRequests(ctx context.Context, userID string) ([]map[string]interface{}, error) {
    var allRequests []map[string]interface{}

    // Construct the document reference for Firestore
    userRef := s.FirestoreClient.Doc("user/" + userID)
    log.Printf("Searching for ads with ParentId: %v\n", userRef)

    // Step 1: Get all ads for the specified user (ParentId as a DocumentRef)
    adsIter := s.FirestoreClient.Collection("ad").Where("ParentId", "==", userRef).Documents(ctx)

    foundAds := false

    for {
        adDoc, err := adsIter.Next()
        if err == iterator.Done {
            break
        }
        if err != nil {
            log.Printf("Error retrieving ads: %v\n", err)
            return nil, err
        }

        foundAds = true
        log.Printf("Ad document found with ID: %s\n", adDoc.Ref.ID)

        // Step 2: Retrieve and process the requests field
        requestsField, exists := adDoc.Data()["requests"]
        if !exists {
            log.Printf("No requests field found for ad with ID: %s\n", adDoc.Ref.ID)
            continue // Skip this ad if there's no requests field
        }

        // Ensure requestsField is of type []interface{} to cast each element
        requestsArray, ok := requestsField.([]interface{})
        if !ok {
            log.Printf("requests field is not of type []interface{} for ad with ID: %s\n", adDoc.Ref.ID)
            continue
        }

        log.Printf("Number of request references found for ad %s: %d\n", adDoc.Ref.ID, len(requestsArray))

        // Process each request reference
        for _, req := range requestsArray {
            reqRef, ok := req.(*firestore.DocumentRef)
            if !ok {
                log.Printf("Invalid request reference in ad %s\n", adDoc.Ref.ID)
                continue
            }

            reqDoc, err := reqRef.Get(ctx)
            if err != nil {
                log.Printf("Error retrieving request document: %v\n", err)
                return nil, err
            }
            allRequests = append(allRequests, reqDoc.Data())
            log.Printf("Request document added with data: %v\n", reqDoc.Data())
        }
    }

    if !foundAds {
        log.Println("No ads found with the specified ParentId.")
    }

    log.Printf("Total requests found: %d\n", len(allRequests))
    return allRequests, nil
}




