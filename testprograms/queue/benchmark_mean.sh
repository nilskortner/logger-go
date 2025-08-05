#!/bin/bash

# Check if a program was provided
if [ -z "$1" ]; then
    echo "Usage: $0 '<command_to_run>'"
    echo "Example: $0 'java -jar 10000_virtualthreads_sleep.jar'"
    exit 1
fi

PROGRAM="$1"
RUNS=10

# Accumulators
total_mrss=0
total_vcs=0
total_ivcs=0

echo "Running benchmark $RUNS times for: $PROGRAM"
echo

for i in $(seq 1 $RUNS); do
    echo "Run #$i..."
    output=$( /usr/bin/time -v bash -c "$PROGRAM" bash 2>&1 )

    # Extract metrics
    mrss=$(echo "$output" | grep -i "Maximum resident set size" | head -n1 | awk -F: '{print $2}' | tr -dc '0-9')
    vcs=$(echo "$output" | grep -i "Voluntary context switches" | head -n1 | awk -F: '{print $2}' | tr -dc '0-9')
    ivcs=$(echo "$output" | grep -i "Involuntary context switches" | head -n1 | awk -F: '{print $2}' | tr -dc '0-9')




    echo "  MRSS: $mrss KB | VCS: $vcs | IVCS: $ivcs"

    total_mrss=$((total_mrss + mrss))
    total_vcs=$((total_vcs + vcs))
    total_ivcs=$((total_ivcs + ivcs))
done

# Calculate averages
mean_mrss=$((total_mrss / RUNS))
mean_vcs=$((total_vcs / RUNS))
mean_ivcs=$((total_ivcs / RUNS))

echo
echo "==== AVERAGE RESULTS ===="
echo "Average MRSS : $mean_mrss KB"
echo "Average VCS  : $mean_vcs"
echo "Average IVCS : $mean_ivcs"
