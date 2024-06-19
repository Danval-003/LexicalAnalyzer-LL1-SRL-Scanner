import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SimulaterComponent } from './simulater.component';

describe('SimulaterComponent', () => {
  let component: SimulaterComponent;
  let fixture: ComponentFixture<SimulaterComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [SimulaterComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(SimulaterComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
