//
//  WHOMessageCell.h
//  Hooli
//
//  Created by dylan on 4/19/14.
//  Copyright (c) 2014 whoisdylan. All rights reserved.
//

#import <UIKit/UIKit.h>

@interface WHOMessageCell : UITableViewCell
@property (strong, nonatomic) IBOutlet UILabel *messageLabel;
@property (strong, nonatomic) IBOutlet UILabel *authorLabel;
@property (strong, nonatomic) IBOutlet UILabel *distanceLabel;

@end
