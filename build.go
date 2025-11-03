package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/a-h/templ"

	"campaign/data"
	"campaign/templates"
)

func generateStaticSite() error {
	fmt.Println("🔨 Generating static site...")

	os.MkdirAll("public", 0755)
	os.MkdirAll("public/blog", 0755)
	os.MkdirAll("public/js", 0755)

	// Copy vendor files
	if err := copyFile("static-vendor/htmx.min.js", "public/js/htmx.min.js"); err != nil {
		return fmt.Errorf("failed to copy htmx: %w", err)
	}
	fmt.Println("✅ Copied: static-vendor/htmx.min.js -> public/js/htmx.min.js")

	if err := renderToFile("public/index.html", templates.Home()); err != nil {
		return err
	}
	fmt.Println("✅ Generated: public/index.html")

	if err := renderToFile("public/about.html", templates.About()); err != nil {
		return err
	}
	fmt.Println("✅ Generated: public/about.html")

	posts := data.GetBlogPosts()
	if err := renderToFile("public/blog.html", templates.BlogList(posts)); err != nil {
		return err
	}
	fmt.Println("✅ Generated: public/blog.html")

	for _, post := range posts {
		filename := filepath.Join("public/blog", post.Slug+".html")
		if err := renderToFile(filename, templates.BlogPost(post)); err != nil {
			return err
		}
		fmt.Printf("✅ Generated: %s\n", filename)
	}

	if err := os.WriteFile("public/style.css", []byte("/* Your styles here */\n"), 0644); err != nil {
		return err
	}
	fmt.Println("✅ Created: public/style.css")

	fmt.Println("✅ Static site generation complete!")
	return nil
}

func renderToFile(filename string, component templ.Component) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return component.Render(context.Background(), file)
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
