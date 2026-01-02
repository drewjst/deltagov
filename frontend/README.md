# Frontend

This project was generated using [Angular CLI](https://github.com/angular/angular-cli) version 21.0.4.

## State Management

This project uses **NGRX SignalStore** (`@ngrx/signals`) for reactive state management. SignalStore provides:

- Signal-based reactivity integrated with Angular's change detection
- Type-safe state with `withState()`
- Derived state via `withComputed()`
- Actions/methods via `withMethods()`

### Store Pattern

Feature stores are co-located with their components:

```
src/app/components/living-bill/
├── living-bill.ts
├── living-bill.html
├── living-bill.scss
└── living-bill.store.ts    # SignalStore for this feature
```

### Usage

```typescript
// Inject the store
readonly store = inject(LivingBillStore);

// Read state (signals)
store.versions()
store.isLoading()

// Call methods
store.selectFromVersion('v1');
store.setError('Something went wrong');
```

## Development server

To start a local development server, run:

```bash
ng serve
```

Once the server is running, open your browser and navigate to `http://localhost:4200/`. The application will automatically reload whenever you modify any of the source files.

## Code scaffolding

Angular CLI includes powerful code scaffolding tools. To generate a new component, run:

```bash
ng generate component component-name
```

For a complete list of available schematics (such as `components`, `directives`, or `pipes`), run:

```bash
ng generate --help
```

## Building

To build the project run:

```bash
ng build
```

This will compile your project and store the build artifacts in the `dist/` directory. By default, the production build optimizes your application for performance and speed.

## Running unit tests

To execute unit tests with the [Vitest](https://vitest.dev/) test runner, use the following command:

```bash
ng test
```

## Running end-to-end tests

For end-to-end (e2e) testing, run:

```bash
ng e2e
```

Angular CLI does not come with an end-to-end testing framework by default. You can choose one that suits your needs.

## Additional Resources

For more information on using the Angular CLI, including detailed command references, visit the [Angular CLI Overview and Command Reference](https://angular.dev/tools/cli) page.
