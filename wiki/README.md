# SupaManager Wiki Pages

This directory contains all the markdown files for the SupaManager GitHub Wiki.

---

## Publishing to GitHub Wiki

GitHub Wikis are actually Git repositories themselves. Here's how to publish these pages:

### Method 1: Clone Wiki Repository (Recommended)

```bash
# 1. Go to your GitHub repo wiki page
# https://github.com/YOUR_USERNAME/supabase-manager/wiki

# 2. Clone the wiki repository
git clone https://github.com/YOUR_USERNAME/supabase-manager.wiki.git

# 3. Copy all markdown files from this directory
cp /home/haider/supabase-manager/wiki/*.md supabase-manager.wiki/

# 4. Commit and push
cd supabase-manager.wiki
git add .
git commit -m "Add comprehensive wiki documentation"
git push origin master
```

### Method 2: Manual Upload via GitHub Interface

1. Go to your repository on GitHub
2. Click the "Wiki" tab
3. For each markdown file:
   - Click "New Page"
   - Use the filename (without .md) as the page title
   - Paste the content
   - Click "Save Page"

---

## Wiki Pages

### Core Pages

1. **Home.md** - Landing page with navigation
2. **Quick-Start.md** - 5-minute setup guide
3. **Architecture-Overview.md** - System architecture
4. **API-Reference.md** - Complete API documentation
5. **Configuration-Reference.md** - All config options
6. **Troubleshooting.md** - Problem solutions
7. **FAQ.md** - Frequently asked questions
8. **Supabase-Stack.md** - 12-service architecture
9. **Roadmap.md** - Development plan

### Sidebar

- **_Sidebar.md** - Navigation sidebar for GitHub Wiki

### Additional Pages to Create

These pages are referenced but not yet created:

- Installation-Guide.md
- First-Steps.md
- Environment-Variables.md
- Docker-Networking.md
- Service-Communication.md
- Database-Schema.md
- Development-Guide.md
- Code-Structure.md
- Running-Tests.md
- Contributing.md
- Deployment.md
- Monitoring.md
- Backup-Recovery.md
- CLI-Commands.md
- Changelog.md

---

## Page Linking

### Internal Links

Use wiki-style links (no .md extension):

```markdown
See the [Architecture Overview](Architecture-Overview) for details.
```

### External Links

```markdown
[Supabase Docs](https://supabase.com/docs)
```

### Anchor Links

```markdown
See [Authentication](#authentication) below.
```

---

## Formatting Guidelines

### Code Blocks

Use triple backticks with language:

    ```bash
    docker compose up -d
    ```

### Admonitions

GitHub Wiki supports blockquotes for notes:

```markdown
> **Note:** This is an important note.

> **Warning:** This is a warning.
```

### Tables

```markdown
| Column 1 | Column 2 |
|----------|----------|
| Data 1   | Data 2   |
```

### Sections

Use `---` for horizontal rules between major sections.

---

## Maintenance

### Updating Pages

```bash
# 1. Clone wiki repo
git clone https://github.com/YOUR_USERNAME/supabase-manager.wiki.git

# 2. Edit files
cd supabase-manager.wiki
vim Home.md

# 3. Commit and push
git add Home.md
git commit -m "Update Home page"
git push origin master
```

### Adding New Pages

1. Create the .md file in this directory
2. Update _Sidebar.md with link to new page
3. Copy to wiki repository
4. Commit and push

### Versioning

Consider adding version numbers to pages that change frequently:

```markdown
**Last Updated:** November 2025
**Version:** 1.0
```

---

## Best Practices

### 1. Keep Pages Focused

Each page should cover one main topic.

### 2. Use Consistent Structure

- Title at top
- Overview section
- Detailed sections
- Related links at bottom

### 3. Cross-Reference

Link between related pages:
```markdown
See also: [API Reference](API-Reference)
```

### 4. Update Regularly

- Keep documentation in sync with code
- Update when features change
- Fix outdated information

### 5. Test Links

After publishing, check all links work.

---

## Wiki Structure

```
wiki/
├── Home.md                      # Landing page
├── _Sidebar.md                  # Navigation menu
│
├── Getting Started/
│   ├── Quick-Start.md
│   ├── Installation-Guide.md
│   └── First-Steps.md
│
├── Configuration/
│   ├── Configuration-Reference.md
│   ├── Environment-Variables.md
│   └── Docker-Networking.md
│
├── Architecture/
│   ├── Architecture-Overview.md
│   ├── Supabase-Stack.md
│   ├── Service-Communication.md
│   └── Database-Schema.md
│
├── Development/
│   ├── API-Reference.md
│   ├── Development-Guide.md
│   ├── Code-Structure.md
│   ├── Running-Tests.md
│   └── Contributing.md
│
├── Operations/
│   ├── Deployment.md
│   ├── Monitoring.md
│   ├── Backup-Recovery.md
│   └── Troubleshooting.md
│
└── Reference/
    ├── CLI-Commands.md
    ├── FAQ.md
    ├── Roadmap.md
    └── Changelog.md
```

---

## Current Status

### ✅ Created

- Home.md
- Quick-Start.md
- Architecture-Overview.md
- API-Reference.md
- Configuration-Reference.md
- Troubleshooting.md
- FAQ.md
- Supabase-Stack.md
- Roadmap.md
- _Sidebar.md

### 📝 To Create

- Installation-Guide.md (can copy from README)
- First-Steps.md (tutorial for first project)
- Environment-Variables.md (detailed env var reference)
- Docker-Networking.md (container networking guide)
- Service-Communication.md (how services talk)
- Database-Schema.md (schema reference)
- Development-Guide.md (dev environment setup)
- Code-Structure.md (codebase walkthrough)
- Running-Tests.md (testing guide)
- Contributing.md (contribution guidelines)
- Deployment.md (production deployment)
- Monitoring.md (monitoring setup)
- Backup-Recovery.md (backup strategies)
- CLI-Commands.md (command reference)
- Changelog.md (version history)

---

## Publishing Checklist

Before publishing to GitHub Wiki:

- [ ] All pages created
- [ ] Links tested (internal and external)
- [ ] Code examples tested
- [ ] Screenshots added (if needed)
- [ ] Sidebar navigation updated
- [ ] Cross-references added
- [ ] Spelling checked
- [ ] Consistent formatting
- [ ] Version numbers added
- [ ] Last updated dates added

---

## Need Help?

- GitHub Wiki documentation: https://docs.github.com/en/communities/documenting-your-project-with-wikis
- Markdown guide: https://guides.github.com/features/mastering-markdown/
- Ask in Discord: https://discord.gg/4k5HRe6YEp

---

**Ready to publish?** Follow the instructions in "Publishing to GitHub Wiki" above!
