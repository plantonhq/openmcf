# DynamoDB Stream to Lambda

Change-data capture without an event-source mapping: the pipe reads the table's stream from LATEST, bisects failing batches automatically, dead-letters what still fails, and invokes the audit function synchronously so failures retry. The starting position is fixed for life — pick it deliberately.
