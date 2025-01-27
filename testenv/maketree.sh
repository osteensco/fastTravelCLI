#!/bin/bash

maketree() {

    if [[ -d "testdirtree" ]]; then
        echo "Project tree already exists, building navigation history..."

        pushd testdirtree
        pushd branch_1_a
        pushd branch_1_b 
        pushd branch_1_c
        pushd ../../../
        pushd branch_2_a
        pushd branch_2_b
        pushd branch_2_c
        pushd ../../../
        pushd branch_3_a 
        pushd branch_3_b 
        pushd branch_3_c
        pushd ../../../../

        return    
    fi

    echo "Building out project tree and navigation history...."
    mkdir testdirtree
    pushd testdirtree
    mkdir branch_1_a
    pushd branch_1_a
    mkdir branch_1_b
    pushd branch_1_b 
    mkdir branch_1_c
    pushd branch_1_c
    pushd ../../../
    mkdir branch_2_a
    pushd branch_2_a
    mkdir branch_2_b
    pushd branch_2_b
    mkdir branch_2_c
    pushd branch_2_c
    pushd ../../../
    mkdir branch_3_a 
    pushd branch_3_a 
    mkdir branch_3_b 
    pushd branch_3_b 
    mkdir branch_3_c 
    pushd branch_3_c
    pushd ../../../../
    
    tree .
}

# maketree creates a directory history in which to test the available_commands off of.
maketree
