---
name: erd-maker
description: Use when asked to generate an Entity-Relationship Diagram (ERD). Generates ERD using Mermaid scripting language.
---

# ERD Maker 

Generate Entity-Relationship Diagrams (ERD) using Mermaid scripting language. This skill analyzes SQL migration files to extract database schema information and produces a visual representation of the database structure in an `ERD.md` file.

## Instructions

### Step 1: Identify Database Schema

- Analyze all .sql files in the `migrations/` directory to understand the database schema.
- Extract table definitions, including columns, data types, primary keys, foreign keys, and relationships between tables.

### Step 2: Generate Mermaid Script

- Create a Mermaid script that represents the ERD based on the extracted schema information.
- Use the following Mermaid syntax for tables and relationships:

```mermaid
erDiagram
    TableName {
        ColumnName1 DataType
        ColumnName2 DataType
        ...
    }
    TableName1 ||--o{ TableName2 : "relationship"
```

### Step 3: Generate ERD.md

- Create a file named `ERD.md` in the root directory of the project.
- Write the generated Mermaid script into `ERD.md` so that it can be rendered by Mermaid-compatible tools.