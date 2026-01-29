# Aliasly Product Landing Page

## Project Context

This is the static landing page for **Aliasly** - a cross-platform CLI tool for command alias management.

**Main Repository:** https://github.com/Eganathan/aliasly

## File Structure

```
product-landing-page/
├── index.html      # Main landing page
├── styles.css      # Stylesheet
├── CLAUDE.md       # This context file
└── assets/         # Images and icons (if any)
```

## Key Product Features to Highlight

1. **No Prefix Required** - Aliases run directly (e.g., `gs` instead of `al gs`)
2. **Parameterized Commands** - Dynamic values with `{{param}}` syntax
3. **Web UI Configuration** - Visual interface via `al config`
4. **Import/Export** - YAML backup and restore
5. **Cross-Platform** - macOS and Linux support
6. **Single Binary** - No dependencies

## Important Links

- GitHub: https://github.com/Eganathan/aliasly
- Install Script: https://raw.githubusercontent.com/Eganathan/aliasly/refs/heads/master/scripts/install.sh

## CLI Commands Reference

| Command | Purpose |
|---------|---------|
| `al list` | Display all aliases |
| `al add` | Create new alias interactively |
| `al remove <name>` | Delete an alias |
| `al config` | Launch web configuration UI |
| `al export [file]` | Export configuration |
| `al import <file>` | Import configuration |

## Update Guidelines

- Keep the design minimal and clean
- Ensure install commands are up-to-date with main repo
- Feature list should match GitHub README
- Test all external links before deployment

## Tech Stack

- Pure HTML/CSS (no JavaScript)
- No build step required
- Hosted on GitHub Pages

## GitHub Pages Setup

1. Go to repository Settings > Pages
2. Set source to "Deploy from a branch"
3. Select the branch and folder containing these files
4. Site will be available at `https://eganathan.github.io/aliasly/` (or custom domain)
