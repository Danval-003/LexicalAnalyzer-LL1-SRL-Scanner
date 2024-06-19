import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TableSlrComponent } from './table-slr.component';

describe('TableSlrComponent', () => {
  let component: TableSlrComponent;
  let fixture: ComponentFixture<TableSlrComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TableSlrComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(TableSlrComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
