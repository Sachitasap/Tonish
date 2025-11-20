# MyFlowBook - Quick Start Guide

## 🎯 What You Can Do Now

### 1. Create and Organize Notebooks 📓
```
MyFlowBook Page → Click "+ New Notebook"
├── Enter notebook name (e.g., "Work Notes", "Personal Journal")
├── Add tags (e.g., "work, projects, important")
└── Click "Create Notebook"
```

**Features**:
- Unlimited notebooks
- Search notebooks by name or tags
- Pin important notebooks (shows at top with special styling)
- See page count for each notebook

---

### 2. Create and Edit Pages ✍️
```
Open Notebook → Click "+ New Page"
├── Enter page title
├── Write content with markdown support
│   ├── **bold text**
│   ├── *italic text*
│   ├── `code snippets`
│   └── # Headings
├── Add tags
└── Click "Create Page"
```

**Editing**:
- Click page from list to view
- Click "✏️ Edit" to modify
- Use toolbar buttons for quick formatting
- Click "💾 Save" when done
- Toggle between Edit and Preview modes

---

### 3. Search Your Notes 🔍
```
Notebook Page → Type in search bar → Press Enter
├── Searches titles and content
├── Shows matching pages
└── Click result to open page
```

**Search Locations**:
- Main page: Search notebooks by name/tags (instant filter)
- Notebook page: Full-text search across all pages (backend-powered)

---

### 4. Pin Important Items 📌
```
Pin Notebooks:
├── Click 📍 icon on any notebook → Becomes 📌
├── Pinned notebooks appear in separate section at top
└── Click 📌 to unpin → Returns to regular section

Pin Pages:
├── Click 📍 icon on any page → Becomes 📌
├── Pinned pages appear first in list
└── Yellow border marks pinned pages
```

---

### 5. Export Your Work 📤
```
Open Notebook → Click "📤 Export" dropdown
├── 📄 Export as Text
│   └── Clean, readable .txt file
├── 📊 Export as CSV
│   └── Spreadsheet-compatible format
└── 📝 Export as Markdown
    └── GitHub/GitLab ready .md file
```

**What Gets Exported**:
- All pages in the notebook
- Titles and full content
- Tags
- Created and updated timestamps
- Formatted for easy reading

---

## 🎨 UI Overview

### Main MyFlowBook Page
```
┌─────────────────────────────────────────────────────┐
│ MyFlowBook 📓              [+ New Notebook]         │
├─────────────────────────────────────────────────────┤
│ 🔍 Search notebooks by name or tags...              │
├─────────────────────────────────────────────────────┤
│ 📌 Pinned Notebooks (2)                             │
│ ┌──────────────┐ ┌──────────────┐                  │
│ │ Work Notes  │ │ Daily Log   │                    │
│ │ 📄 12 pages  │ │ 📄 5 pages   │                  │
│ │ #work #proj  │ │ #personal    │                  │
│ │ [Open][🗑️]   │ │ [Open][🗑️]   │                  │
│ └──────────────┘ └──────────────┘                  │
├─────────────────────────────────────────────────────┤
│ All Notebooks (5)                                   │
│ ┌──────────────┐ ┌──────────────┐ ┌──────────────┐│
│ │ Ideas       │ │ Recipes     │ │ Research    ││
│ │ ...          │ │ ...          │ │ ...          ││
│ └──────────────┘ └──────────────┘ └──────────────┘│
└─────────────────────────────────────────────────────┘
```

### Notebook Detail Page
```
┌─────────────────────────────────────────────────────┐
│ ← Back  Work Notes 📌    [📤 Export][+ New Page]   │
├─────────────────────────────────────────────────────┤
│ 🔍 Search pages...                      [Search]    │
├─────────────────────────────────────────────────────┤
│                                                     │
│ ┌──────────────┬─────────────────────────────────┐ │
│ │ Pages (12)   │ Meeting Notes - 2024-11-15      │ │
│ │              │                                 │ │
│ │ 📌 Meeting   │ # Agenda                        │ │
│ │    Notes     │                                 │ │
│ │    #import   │ 1. Project timeline review      │ │
│ │              │ 2. Budget discussion            │ │
│ │ 📍 Sprint    │ 3. Team updates                 │ │
│ │    Planning  │                                 │ │
│ │    #dev      │ **Action Items:**               │ │
│ │              │ - Follow up with design team    │ │
│ │ 📍 Ideas     │ - Update project plan           │ │
│ │    #brain    │                                 │ │
│ │              │ Tags: #work #meetings #Q4       │ │
│ │              │                                 │ │
│ │              │ [✏️ Edit] [🗑️ Delete]            │ │
│ │              │                                 │ │
│ │              │ Created: Nov 15, 2024           │ │
│ │              │ Updated: Nov 18, 2024           │ │
│ └──────────────┴─────────────────────────────────┘ │
└─────────────────────────────────────────────────────┘
```

### Edit Mode
```
┌─────────────────────────────────────────────────────┐
│ Meeting Notes               [👁️ Preview][🗑️ Delete] │
├─────────────────────────────────────────────────────┤
│ [B][I][</>][H1]                         [💾 Save]   │
├─────────────────────────────────────────────────────┤
│                                                     │
│ # Agenda                                            │
│                                                     │
│ 1. Project timeline review                          │
│ 2. Budget discussion                                │
│ 3. Team updates                                     │
│                                                     │
│ **Action Items:**                                   │
│ - Follow up with design team                        │
│ - Update project plan                               │
│                                                     │
│                                                     │
│ 💡 Tip: Use **bold**, *italic*, `code`, # headings │
└─────────────────────────────────────────────────────┘
```

---

## 🚦 Getting Started (3 Steps)

1. **Start the app**:
   ```bash
   cd /path/to/Tonish
   ./start-dev.sh
   ```

2. **Create your first notebook**:
   - Go to http://localhost:5173/myflowbook
   - Click "+ New Notebook"
   - Name it and add tags
   - Create!

3. **Add your first page**:
   - Open the notebook
   - Click "+ New Page"
   - Write your content
   - Save!

---

## 💡 Pro Tips

### Markdown Quick Reference
```markdown
**bold**                  → bold
*italic*                  → italic
`code`                    → code
# Heading 1              → Heading 1
## Heading 2             → Heading 2
### Heading 3            → Heading 3
```

### Organization Tips
1. Use **tags** liberally - easier to search later
2. **Pin** your most-used notebooks
3. Create **template pages** for recurring content
4. Use **descriptive titles** for easy scanning
5. **Export** regularly as backup

### Search Best Practices
- Search is case-insensitive
- Searches both titles and content
- Use specific terms for better results
- Clear search to see all pages again

---

## 📊 Feature Summary

| Feature | Location | Action |
|---------|----------|--------|
| Create Notebook | Main page | Click "+ New Notebook" |
| Search Notebooks | Main page | Type in search bar |
| Pin Notebook | Notebook card | Click 📍 icon |
| Create Page | Notebook page | Click "+ New Page" |
| Edit Page | Page view | Click "✏️ Edit" |
| Search Pages | Notebook page | Use search bar |
| Pin Page | Page list | Click 📍 on page |
| Export | Notebook page | Click "📤 Export" |

---

## 🎉 You're All Set!

MyFlowBook is ready to help you:
- ✅ Organize your thoughts
- ✅ Keep meeting notes
- ✅ Store ideas
- ✅ Maintain journals
- ✅ Build knowledge base
- ✅ Track projects

**Start writing!** 📝
