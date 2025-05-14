# Scheduler

- [ ] (On startup) Retrieve monitors from db and start the monitor
- [ ] Check for monitors to execute
- [ ] Publish monitor event

# Worker

- [ ] Listen to monitor events
- [ ] Executing the monitor event
- [ ] Publish the results

# Monitor

- [ ] Listen to the Worker events
- [ ] Update database
- [ ] Implement analytics

# Web Service

- [ ] Auth
- [ ] Create, Delete, update monitor

# Flusher

- [ ] Save the redis data to db every x time
