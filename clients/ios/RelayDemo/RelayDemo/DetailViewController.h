//
//  DetailViewController.h
//  RelayDemo
//
//  Created by Sri Panyam on 28/02/2015.
//  Copyright (c) 2015 Panyam. All rights reserved.
//

#import <UIKit/UIKit.h>

@interface DetailViewController : UIViewController

@property (strong, nonatomic) id detailItem;
@property (weak, nonatomic) IBOutlet UILabel *detailDescriptionLabel;

@end

