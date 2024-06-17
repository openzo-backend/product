#!/bin/bash

# Step 1: Backup
git clone https://github.com/openzo-backend/product
cd repository
git branch backup
git push origin backup

# Step 2: Remove History
git checkout --orphan latest_branch
git add -A
git commit -am "Initial commit - Clean history"
git branch -D main
git branch -m main
git push -f origin main

# Step 3: Add .gitignore
cat <<EOL > .gitignore
/config/*.yaml
*.properties
test.db
/vendor/
EOL
git add .gitignore
git commit -m "Add .gitignore"
git push origin main
