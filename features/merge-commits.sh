#!/bin/bash

cd test

# Create branch test
git checkout -b test

# Create several commits on test branch
touch file3
git add file3
git commit -F- <<EOF
feat(file3) : new file 3

create a new file 3
EOF

touch file4
git add file4
git commit -F- <<EOF
feat(file4) : new file 4

create a new file 4
EOF

# Create branch test1
git checkout -b test1

touch file5
git add file5
git commit -F- <<EOF
feat(file5) : new file 5

create a new file 5
EOF

touch file6
git add file6
git commit -F- <<EOF
feat(file6) : new file 6

create a new file 6
EOF

# Create branch test2
git checkout -b test2

touch file7
git add file7
git commit -F- <<EOF
feat(file7) : new file 7

create a new file 7
EOF

touch file8
git add file8
git commit -F- <<EOF
feat(file8) : new file 8

create a new file 8
EOF

# Checkout branch test1
git checkout test1

# Merge branch test2
git merge --no-ff test2

# Checkout branch test
git checkout test

# Merge branch test1
git merge --no-ff test1

# Checkout branch master
git checkout master

# Merge branch test
git merge --no-ff test

# Create several commits on master branch

touch file9
git add file9
git commit -F- <<EOF
feat(file9) : new file 9

create a new file 9
EOF

touch file10
git add file10
git commit -F- <<EOF
feat(file10) : new file 10

create a new file 10
EOF
