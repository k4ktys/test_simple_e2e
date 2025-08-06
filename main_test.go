package main

import (
	"log"
	"testing"

	"github.com/playwright-community/playwright-go"
	"github.com/stretchr/testify/assert"
)

func TestPositive(t *testing.T) {
	testCase := "Message1"

	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start playwright: %v", err)
	}

	browser, err := pw.Chromium.Launch()
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}

	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}

	if _, err = page.Goto("http://localhost:8080"); err != nil {
		log.Fatalf("could not goto: %v", err)
	}

	if err = page.GetByPlaceholder("message").Fill(testCase); err != nil {
		log.Fatalf("could not fill input: %v", err)
	}

	if err = page.GetByText("Send").Click(); err != nil {
		log.Fatalf("could not click on send button: %v", err)
	}

	if _, err = page.Reload(); err != nil {
		log.Fatalf("could not reload page: %v", err)
	}

	entries, err := page.Locator(".list li").All()
	if err != nil {
		log.Fatalf("could not locate: %v", err)
	}

	text, err := entries[0].TextContent()
	if err != nil {
		log.Fatalf("could not get text %v", err)
	}

	assert.Equal(t, len(entries), 1)
	assert.Equal(t, text, testCase)

	if err = browser.Close(); err != nil {
		log.Fatalf("could not close browser: %v", err)
	}

	if err = pw.Stop(); err != nil {
		log.Fatalf("could not stop Playwright: %v", err)
	}
}
