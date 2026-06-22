# Mentoring Strategy: working with a Junior Developer

This is a classic leadership dilemma. Given that your colleague is a junior developer, I strongly recommend a **"Big Picture Architecture + Small Task Execution"** hybrid approach.

Here is the strategy that usually works best for mentoring juniors:

---

## 1. The "Big Picture" (Do this NOW)
You should complete the overall architecture (which we have mostly done with the `PLAN.md` files).

* **Why:** A junior needs to know *why* they are building something. If she doesn't see how the API connects to MSAL, she might build things in a way that needs to be deleted later.
* **What to show her:** The Mermaid diagram, the OpenAPI contract, and the Folder Structure. This gives her "the rails" to stay on.

---

## 2. The Task List (Do this WEEKLY)
While you should have the milestones mapped out for the whole project (e.g., "Week 4: Auth working"), you should only give her a detailed task list one week at a time.

* **Avoid the "Wall of Work":** Giving a junior an 8-week list of 50 tasks can be paralyzing. They often don't know where to start or how to prioritize.
* **The Weekly Sync:** Every Monday, give her 3 to 5 specific "Definition of Done" tasks.
* **Example:** *"By Friday, I want to see the `JobCard` component rendered with mock data from MSW."*
* **Flexibility:** Juniors often hit roadblocks. If you plan 8 weeks out, and she gets stuck on a CSS issue for 3 days in Week 1, your entire 8-week plan is now "wrong," which can be discouraging for her.

---

## 3. How to Guide Her Effectively
* **The "Definition of Done" (DoD):** For every task, tell her exactly what success looks like (e.g., *"The button must turn Red on hover and be accessible via keyboard."*).
* **Code Reviews as Mentorship:** Don't just fix her code. Leave comments explaining why a different approach might be better. This is where she will learn the most from you.
* **The 30-Minute Rule:** Tell her: *"If you are stuck on a problem for more than 30 minutes, stop and ask me."* This prevents her from spinning her wheels for a whole day on a small bug.

---

## Suggested Next Steps
Write down 4 to 5 Major Milestones:
1. Mock UI Setup
2. Auth Integration
3. Real Data Connection
4. Polishing and Feedback

Then, just focus on writing the detailed, bite-sized tasks for Week 1.
