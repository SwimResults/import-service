package importer

import (
	"context"
	"fmt"
	"github.com/openai/openai-go" // imported as openai
	"github.com/openai/openai-go/option"
	athleteClient "github.com/swimresults/athlete-service/client"
	"github.com/swimresults/athlete-service/model"
	importModel "github.com/swimresults/import-service/model"
	"github.com/swimresults/service-core/misc"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func ImportCertificates(directory string, meeting string) (*importModel.ImportCertificateStats, error) {
	pwd, _ := exec.Command("pwd").Output()
	fmt.Printf("reading certificates from: %s\n", string(pwd))

	systemPath := "/app/assets/files/"

	osDir := filepath.Join(systemPath, directory)

	entries, err := os.ReadDir(osDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		path := osDir + e.Name()
		filename := substr(e.Name(), ".")
		outDir := osDir + filename + "/"
		outfile := outDir + e.Name()

		amount, err := exec.Command("qpdf", "--show-npages", path).Output()
		if err != nil {
			println("fatal")
			log.Fatal(err)
		}
		as := string(amount)
		as = strings.Replace(as, "\n", "", -1)
		n, _ := strconv.Atoi(as)

		fmt.Printf("%s (%d)\n", e.Name(), n)

		err = os.MkdirAll(outDir, 0755)
		if err != nil {
			log.Fatal(err)
		}
		err = exec.Command("qpdf", "--split-pages", path, outfile).Run()
		if err != nil {
			println("split pages failed")
			log.Fatalf("Command finished with error: %v\n", err)
		}
	}

	// +-------------------------------+
	// | GET ALL ATHLETES FOR MEETING  |
	// +-------------------------------+

	var ac = athleteClient.NewAthleteClient("https://api.swimresults.de/athlete/v1/")
	athletes, err := ac.GetAthletesByMeeting(meeting)
	if err != nil {
		log.Fatal(err)
	}

	// +-------------------------------+
	// |    INIT OPEN AI CONNECTION    |
	// +-------------------------------+

	client := openai.NewClient(
		option.WithAPIKey(os.Getenv("SR_IMPORT_OPENAI_KEY")),
	)

	// +-------------------------------+
	// | READ SINGLE FILE CERTIFICATES |
	// +-------------------------------+

	dirs, err := os.ReadDir(osDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, d := range dirs {
		if !d.IsDir() {
			continue
		}

		files, err := os.ReadDir(osDir + d.Name() + "/")
		if err != nil {
			println("fatal")
			log.Fatal(err)
		}

		for _, f := range files {
			pdf, err := ReadPdf(filepath.Join(osDir, d.Name(), f.Name()))
			if err != nil {
				println("could not read pdf")
				log.Fatal(err)
			}

			println(pdf)
			pdfA := misc.Aliasify(pdf)
			println("----=== search for athlete ===----")

			var foundAthlete model.Athlete
			for _, a := range athletes {
				for _, alias := range a.Alias {
					if strings.Contains(pdfA, alias) {
						fmt.Printf("FOUND ATHLETE: %s - %s\n", a.Name, a.Team.Name)
						foundAthlete = a
						break
					}
				}
				if !foundAthlete.Identifier.IsZero() {
					break
				}
			}

			if foundAthlete.Identifier.IsZero() {
				println("\t- could not find athlete")
				continue
			}

			chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
				Messages: []openai.ChatCompletionMessageParamUnion{
					openai.SystemMessage("Aus Urkundentext eines Schwimmwettkampfs einen Dateinamen erzeugen: Strecke + Stil + relevante Zusätze (z.B. Vorlauf, Finale, Masters, Punktbeste Leistung), ohne Name, Verein, Veranstaltung, Platz, Jahrgang, Ort, Datum. Zusätze nur, wenn es klare infos dazu gibt! Ausgabe nur der Name!"),
					openai.UserMessage(pdf),
				},
				Model: openai.ChatModelGPT4_1Mini,
			})
			if err != nil {
				panic(err.Error())
			}

			certName := chatCompletion.Choices[0].Message.Content

			_, _, err = ac.ImportCertificate(certName, foundAthlete.Identifier, meeting, filepath.Join(directory, d.Name(), f.Name()))
			if err != nil {
				println("cert import failed")
				log.Fatal(err)
			} else {
				println("imported!")
			}
		}
	}

	return &importModel.ImportCertificateStats{Amount: 0}, nil
}
