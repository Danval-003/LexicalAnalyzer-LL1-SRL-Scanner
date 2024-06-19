import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { PrincipalComponent } from './pages/principal/principal.component';
import { AboutUsComponent } from './pages/about-us/about-us.component'
import { SimulaterComponent } from './pages/simulater/simulater.component';

export const routes: Routes = [
  // Put principal component into /
  {path: '', component: PrincipalComponent},
  {path: 'about-us', component: AboutUsComponent},
  {path: 'simulator', component: SimulaterComponent}
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})

export class AppRoutingModule { }
