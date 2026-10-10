package bootstrap

import (
	"context"
	"encoding/json"
	"fmt"
	app "github.com/Maaku050/elabtrack-v2/backend/internal/application/borrowing"
	b "github.com/Maaku050/elabtrack-v2/backend/internal/domain/borrowing"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/inventory"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/persistence/postgres"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"os"
	"testing"
	"time"
)

func TestPhase11BrowserServer(t *testing.T) {
	if os.Getenv("ELABTRACK_PHASE11_BROWSER") != "1" {
		t.Skip("explicit isolated browser gate")
	}
	infra := batchInfra(t)
	c := buildContainer(infra)
	users, passwords, _ := phase7Fixtures(t, infra, c)
	ctx := context.Background()
	marker := "TEST P11 " + uuid.NewString()[:8]
	active := true
	category, e := c.InventorySvc.Category(ctx, users["staff"].ID, uuid.Nil, uuid.NewString(), inventory.CategoryInput{Name: marker + " Utensils", Active: &active})
	if e != nil {
		t.Fatal("category fixture")
	}
	equipment := []inventory.Equipment{}
	names := []string{"Wooden Spoon", "Culinary Tongs", "Chafing Dish", "Food Pan", "Muffin Tray", "Blender"}
	files := []string{"wooden-spoon.png", "tongs.png", "chafing-dish.png", "food-pan.png", "muffin-tray.png", "blender.png"}
	for i := 0; i < 26; i++ {
		name := fmt.Sprintf("%s %02d %s", marker, i, names[i%6])
		v, e := c.InventorySvc.Create(ctx, users["staff"].ID, uuid.NewString(), inventory.Create{Metadata: inventory.Metadata{Name: name, Description: "Fictional selection browser fixture", CategoryID: &category.ID}, Opening: 50, Reason: "TEST browser stock"})
		if e != nil {
			t.Fatal("equipment fixture")
		}
		if i < 6 {
			raw, e := os.ReadFile("../../../frontend/src/features/visual-preview/assets/" + files[i])
			if e != nil {
				t.Fatal("local catalog asset")
			}
			v, e = c.InventorySvc.SaveImage(ctx, users["staff"].ID, v.ID, uuid.NewString(), v.Version, raw)
			if e != nil {
				t.Fatal("validated protected catalog image")
			}
		}
		equipment = append(equipment, v)
	}
	past := time.Now().Add(-49 * time.Hour).UTC().Truncate(time.Microsecond)
	due := past.Add(time.Hour)
	svc := app.NewService(phase8Clock{postgres.NewBorrowingRepository(infra.DB.Pool), past}, postgres.NewAccountsRepository(infra.DB.Pool), postgres.NewInventoryRepository(infra.DB.Pool), c.TermsSvc, infra.Tx)
	loan, e := svc.Direct(ctx, users["staff"].ID, uuid.NewString(), b.Input{BorrowerID: users["borrower"].ID, Items: []b.Line{{EquipmentID: equipment[5].ID, Quantity: 3}}, DueAt: &due, Confirm: true, Handover: true})
	if e != nil {
		t.Fatal("own overdue detail fixture")
	}
	fixtures := map[string]any{"equipment": equipment, "loan": loan, "marker": marker, "category": category}
	for role, u := range users {
		fixtures[role] = map[string]any{"email": u.Email, "password": passwords[role]}
	}
	f, e := os.OpenFile("/tmp/elabtrack-phase11-fixtures.json", os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if e != nil {
		t.Fatal("private fixture exists")
	}
	raw, _ := json.Marshal(fixtures)
	_, e = f.Write(raw)
	f.Close()
	if e != nil {
		t.Fatal(e)
	}
	server := RegisterRoutes(newServer(infra.Config, c, infra.Logger), c)
	go func() { _ = server.Listen("127.0.0.1:18085", fiber.ListenConfig{DisableStartupMessage: true}) }()
	defer server.Shutdown()
	for deadline := time.Now().Add(40 * time.Minute); time.Now().Before(deadline); {
		if _, e := os.Stat("/tmp/elabtrack-phase11-browser-stop"); e == nil {
			return
		}
		time.Sleep(200 * time.Millisecond)
	}
	t.Fatal("browser gate timeout")
}
