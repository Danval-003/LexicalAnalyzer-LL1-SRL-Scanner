import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-animated-title',
  standalone: true,
  imports: [],
  templateUrl: './animated-title.component.html',
  styleUrl: './animated-title.component.sass'
})
export class AnimatedTitleComponent {
  @Input() text: string = 'Interstellar';
}
