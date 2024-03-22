package srv

import (
	"fmt"
	"github.com/heather92115/translator/internal/db"
	"github.com/heather92115/translator/internal/mdl"
	"github.com/joho/godotenv"
	"log"
	"math"
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {

	// find the .env.test file
	envPath := os.Getenv("ENV_TEST_PATH")

	// Load the .env.test file
	if err := godotenv.Load(envPath); err != nil {
		log.Fatalf("Error loading .env.test file: %v", err)
	}

	dbUser := os.Getenv("DATABASE_USER")
	dbIp := os.Getenv("DATABASE_IP")
	dbPort := os.Getenv("DATABASE_PORT")
	if dbUser != "tester" {
		log.Fatalf("Did not find the tester user: %s, %s, %s", dbUser, dbIp, dbPort)
	}

	err := db.CreatePool()
	if err != nil {
		fmt.Printf("Failed DB connections, %v\n", err)
		return
	}

	// Now that the environment variables are set, run the tests
	code := m.Run()

	// Exit with the status code from the test run
	os.Exit(code)
}

func TestIntegrationFixitService_CreateFindFixitByID(t *testing.T) {
	// Create an instance of SQL FixitService
	fixitService, err := NewFixitService()
	if err != nil {
		t.Errorf("Unexpected error: %v, failed to create Fixit Service", err)
	}

	testFixit := &mdl.Fixit{
		Status:    "pending",
		FieldName: "Definition",
		Comments:  "Initial comment",
		CreatedBy: "tester",
		Created:   time.Now(),
	}
	_ = fixitService.CreateFixit(testFixit)

	fixit, err := fixitService.FindFixitByID(testFixit.ID)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if fixit.ID != testFixit.ID {
		t.Errorf("Expected fixit ID %d, got %d", testFixit.ID, fixit.ID)
	}

	_, err = fixitService.FindFixitByID(9999999)
	if err == nil {
		t.Error("Expected an error for non-existing fixit, but got nil")
	}
}

// TestIntegrationFixitService_FindFixits tests the functionality of the Fixit Service
func TestIntegrationFixitService_CreateFindUpdate(t *testing.T) {
	// Create an instance of SQL FixitService
	fixitService, err := NewFixitService()
	if err != nil {
		t.Errorf("Unexpected error: %v, failed to create Fixit Service", err)
	}

	testFixit := &mdl.Fixit{
		Status:    "pending",
		FieldName: "Definition",
		Comments:  "Initial comment",
		CreatedBy: "tester",
		Created:   time.Now(),
	}
	err = fixitService.CreateFixit(testFixit)
	if err != nil {
		t.Errorf("Unexpected error on create: %v", err)
	}

	fixitList, err := fixitService.FindFixits("pending", 0, nil, 5)
	if err != nil {
		t.Errorf("Unexpected error on query: %v", err)
	}

	if len(*fixitList) == 0 {
		t.Errorf("Expected to find fixits")
	}

	for _, fixit := range *fixitList {

		if fixit.Status != "pending" {
			t.Errorf("Expected fixit with id %d to have status 'pending'", fixit.ID)
		}

		fixit.Status = "completed"
		updated, err := fixitService.UpdateFixit(&fixit)
		if err != nil {
			t.Errorf("Unexpected error on update: %v", err)
		}
		if updated.Status != "completed" {
			t.Errorf("Expected fixit with id %d to have status 'completed'", fixit.ID)
		}
	}
}

// TestIntegrationVocabService_CreateFindUpdate tests the functionality the Vocab Service
func TestIntegrationVocabService_CreateFindUpdate(t *testing.T) {
	// Create an instance of SQL VocabService
	vocabService, err := NewVocabService()
	if err != nil {
		t.Errorf("Unexpected error: %v, failed to create Vocab Service", err)
	}

	auditService, err := NewAuditService()
	if err != nil {
		t.Errorf("Unexpected error: %v, failed to create Audit Service", err)
	}

	txt := fmt.Sprintf("empecé    %s", randomLetters(10))

	testVocab := &mdl.Vocab{
		// Learning lang must be unique
		LearningLang:     txt,
		FirstLang:        "I began",
		Infinitive:       "empezar",
		Hint:             "past tense",
		Pos:              "verb",
		LearningLangCode: "es",
		KnownLangCode:    "en",
	}

	err = vocabService.CreateVocab(testVocab)
	if err != nil {
		log.Printf("Validation error on create vocab %+v, err: %v", testVocab, err)
		t.Errorf("Unexpected error on create: %v", err)
		return
	}

	currentTime := time.Now()
	twoSecondsAgo := currentTime.Add(-2 * time.Second)
	duration := mdl.Duration{
		Start: twoSecondsAgo,
		End:   currentTime,
	}
	auditList, err := auditService.FindAudits("vocab", &duration, math.MaxInt)
	if err != nil {
		t.Errorf("Unexpected error on audit query: %v", err)
		return
	} else if len(*auditList) == 0 {
		t.Errorf("Expected to find at least one audit record for the vocab table")
	}

	found := false
	for _, audit := range *auditList {
		if audit.ObjectID == testVocab.ID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected to find audit with object id %d, but did not", testVocab.ID)
	}

	vocabList, err := vocabService.FindVocabs("es", true, 5)
	if err != nil {
		t.Errorf("Unexpected error on vocab query: %v", err)
		return
	}

	if len(*vocabList) == 0 {
		t.Errorf("Expected to find vocabs")
		return
	}

	for _, vocab := range *vocabList {

		if vocab.LearningLangCode != "es" {
			t.Errorf("Expected vocab with id %d to have learning lang code 'is'", vocab.ID)
		}

		vocab.Hint = "starts with 'em'"
		updated, err := vocabService.UpdateVocab(&vocab)
		if err != nil {
			t.Errorf("Unexpected error on update %+v, err: %v", vocab, err)
			return
		}
		if updated.Hint != "starts with 'em'" {
			t.Errorf("Expected vocab with id %d to have updated hint", vocab.ID)
		}
	}
}

// This func creates randomness for the learning lang field which must be unique.
func randomLetters(length int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	var result []rune

	for i := 0; i < length; i++ {
		result = append(result, rune(letters[r.Intn(len(letters)-1)]))
	}

	return string(result)
}
