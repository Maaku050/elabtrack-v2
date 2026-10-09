# Inventory filter Chromium evidence

Real Chromium used the existing isolated API on 18085, Vite on 15175 and PostgreSQL on 54832. Generated Admin/Staff/Borrower fixtures and clearly synthetic test-only terms were confined to the disposable database. No owner account or official terms publication was used.

Nine acceptance checks passed. The matrix contained 31 equipment records; 27 were ACTIVE with positive stock and matched the filter. Server pages contained 25 and 2 distinct rows; the matching available-stock aggregate was 46. Screenshots preserve the owner-approved Stage B layout; only the availability label changed.

[Machine-readable acceptance record](acceptance.json) includes sanitized method/path/query evidence, refresh and zero unexpected document navigation checks. No authorization headers, passwords or request bodies are recorded.

- [inactive-availability-empty-light.png](inactive-availability-empty-light.png)
- [inactive-availability-empty-dark.png](inactive-availability-empty-dark.png)
- [availability-320-light.png](availability-320-light.png)
- [availability-320-dark.png](availability-320-dark.png)
- [availability-390-light.png](availability-390-light.png)
- [availability-390-dark.png](availability-390-dark.png)
- [availability-1366-light.png](availability-1366-light.png)
- [availability-1366-dark.png](availability-1366-dark.png)
- [borrower-availability-light.png](borrower-availability-light.png)
- [borrower-availability-dark.png](borrower-availability-dark.png)
