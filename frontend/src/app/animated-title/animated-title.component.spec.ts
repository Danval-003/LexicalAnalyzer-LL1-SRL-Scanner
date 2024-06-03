import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AnimatedTitleComponent } from './animated-title.component';

describe('AnimatedTitleComponent', () => {
  let component: AnimatedTitleComponent;
  let fixture: ComponentFixture<AnimatedTitleComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AnimatedTitleComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(AnimatedTitleComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
