import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { AnimatedTitleComponent } from './animated-title/animated-title.component';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet, AnimatedTitleComponent],
  templateUrl: './app.component.html',
  styleUrl: './app.component.sass'
})
export class AppComponent {
  title = 'frontend';
}
