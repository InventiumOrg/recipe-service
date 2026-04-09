# Business Metrics for Recipe Service

## 🎯 Business-Specific Metrics Added

### Recipe Operations
- **`recipe_operations_total`** - Counter tracking recipe operations
  - Labels: `operation` (get, list, create, update, delete), `recipe_name`

### Database Operations
- **`database_operation_duration_seconds`** - Histogram of DB operation times
  - Labels: `operation`, `table`
  
- **`database_operation_errors_total`** - Counter of DB errors
  - Labels: `operation`, `error_type`

### Authentication & Security
- **`authentication_attempts_total`** - Counter of auth attempts
  - Labels: `operation`, `status`

### Business Logic
- **`recipes_active`** - Gauge of active recipes (best-effort; incremented on create/delete)

## 📊 Example Grafana Queries

### Business KPIs
```promql
# Recipe create rate
sum by (recipe_name) (rate(recipe_operations_total{operation="create"}[5m]))

# Most popular recipe operations
sum by (operation) (rate(recipe_operations_total[5m]))

# Database operation performance
histogram_quantile(0.95, database_operation_duration_seconds)

# Error rates by operation
rate(database_operation_errors_total[5m]) / rate(recipe_operations_total[5m])

# Authentication success rate
(rate(authentication_attempts_total[5m]) - rate(authentication_attempts_total{status="failed"}[5m])) / rate(authentication_attempts_total[5m])
```

### Operational Metrics
```promql
# Active recipes trend
recipes_active

# Database operation errors by type
sum by (error_type) (rate(database_operation_errors_total[5m]))

# Slowest database operations
topk(5, histogram_quantile(0.95, database_operation_duration_seconds) by (operation))
```

## 🚀 Test Your Metrics

1. **Start your service**
2. **Make some API calls:**
   ```bash
   # Health check
   curl http://localhost:7162/healthz
   
   # List recipes
   curl http://localhost:7162/v1/recipe/list
   
   # Create a recipe
   curl -X POST http://localhost:7162/v1/recipe/create \
        -F "name=TestRecipe" -F "ingredients=egg" -F "ingredients=milk" -F "cost=12000"
   ```

3. **Check Grafana Cloud** - Your metrics should appear within 30 seconds

## 🔧 Next Steps

You can extend these metrics by:
- Adding more business operations (update, delete)
- Tracking user-specific metrics
- Adding inventory-related metrics
- Monitoring performance SLAs