# Course Enrollment Technical Documentation

## Endpoints

### POST /academic/course-offering/{id}/enroll

Before the student succeeding the course offering enrollment, we must the validate with these rules:

- No enrollment duplication.
- Check if the registered course offerings is less than course offering capacity. Enrollment will be fail if the registrations for those course offering is fully booked.
- Check for any previously course registration schedule overlaps
  - Each course has a `credit`, each 1 credit is worth 50 minutes. If 2 credits, is 100 minutes and so on
  - Each course offerings has a start time. Expanding `start_time` to `end_time = (start_time + (credit * 50 minutes))` we will get the `start_time` to `end_time` range
  - If the intended enrollment has a schedule overlap to previously enrolled course offerings, the enrollment will be fail.
