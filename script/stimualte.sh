#!/bin/bash

# Configuration
BASE_URL="http://localhost:7162/v1/recipe"

# Number of requests to make (default: 10)
NUM_REQUESTS=${1:-10}

# Delay between requests in seconds (default: 0.5)
DELAY=${2:-0.5}

echo "Starting recipe service simulation with $NUM_REQUESTS requests..."
echo "Delay between requests: ${DELAY}s"
echo "----------------------------------------"

# Function to create a recipe
create_recipe() {
    local name=$1
    local cost=$2
    shift 2
    local ingredients=("$@")
    
    local curl_args=()
    for ing in "${ingredients[@]}"; do
        curl_args+=( --form "ingredients=${ing}" )
    done

    curl --location --request POST "${BASE_URL}/create" \
        --form "name=${name}" \
        --form "cost=${cost}" \
        "${curl_args[@]}" \
        --silent --show-error
}

# Function to get a recipe
get_recipe() {
    local recipe_id=$1
    
    curl --location --request GET "${BASE_URL}/${recipe_id}" \
        --silent --show-error
}

# Function to list recipes
list_recipes() {
    curl --location --request GET "${BASE_URL}/list" \
        --silent --show-error
}

# Function to update a recipe
update_recipe() {
    local recipe_id=$1
    local name=$2
    local cost=$3
    shift 3
    local ingredients=("$@")
    
    local curl_args=()
    for ing in "${ingredients[@]}"; do
        curl_args+=( --form "ingredients=${ing}" )
    done

    curl --location --request PUT "${BASE_URL}/${recipe_id}" \
        --form "name=${name}" \
        --form "cost=${cost}" \
        "${curl_args[@]}" \
        --silent --show-error
}

# Simulate various operations
for i in $(seq 1 $NUM_REQUESTS); do
    echo "Request $i/$NUM_REQUESTS"
    
    # Generate variation in the data
    RECIPE_ID=$((($i % 10) + 1))  # Cycle through recipe IDs 1-10
    NAME="Recipe-${RECIPE_ID}-v${i}"
    COST=$((10000 + ($i * 1500)))
    INGREDIENTS=("ingredient-$((i % 5 + 1))" "ingredient-$((i % 7 + 1))" "ingredient-$((i % 11 + 1))")
    
    # Perform different operations based on request number
    case $((i % 4)) in
        0)
            echo "Creating recipe: name=${NAME}, cost=${COST}, ingredients=${INGREDIENTS[*]}"
            create_recipe "${NAME}" "${COST}" "${INGREDIENTS[@]}"
            ;;
        1)
            echo "Getting recipe with ID: $RECIPE_ID"
            get_recipe $RECIPE_ID
            ;;
        2)
            echo "Listing recipes"
            list_recipes
            ;;
        3)
            echo "Updating recipe ID: $RECIPE_ID with new cost: $COST"
            update_recipe "$RECIPE_ID" "${NAME}" "${COST}" "${INGREDIENTS[@]}"
            ;;
    esac
    
    echo ""
    echo "Response for request $i completed"
    
    # Add delay between requests (except for the last one)
    if [ $i -lt $NUM_REQUESTS ]; then
        echo "Waiting ${DELAY}s before next request..."
        sleep $DELAY
    fi
    
    echo "----------------------------------------"
done

echo ""
echo "Simulation completed!"
echo "Total requests made: $NUM_REQUESTS"