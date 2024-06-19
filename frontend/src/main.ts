import { bootstrapApplication } from '@angular/platform-browser';
import { appConfig } from './app/app.config';
import { AppComponent } from './app/app.component';
import { PrincipalComponent } from './app/pages/principal/principal.component';
import { SimulaterComponent } from './app/pages/simulater/simulater.component';

bootstrapApplication(AppComponent, appConfig)
  .catch((err) => console.error(err));

bootstrapApplication(PrincipalComponent, appConfig)
  .catch((err) => console.error(err));

bootstrapApplication(SimulaterComponent, appConfig)
  .catch((err) => console.error(err));
