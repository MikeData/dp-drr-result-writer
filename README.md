# Dimension Recombination Reporter - Result Writer

WORK IN PROGRESS - 10% project.

NOTE - keeping this as an offline script for now.

The DRR is a microservice based project to analyse all permutations from structural alterations (combinations of dimension items added/removed)
that could be applied to a given source dataset. Presenting the findings as a human readable report.

The result writer checks each result, discarding those that don't meet the filter (for now hardcoded as those that don't reduce sparsity) and writing positive results to CSV before sending back to the controller for further permutations.
